package database

import (
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/gofrs/uuid"
)

func (db *appdbimpl) CreateConversation(currentUserID, targetUsername string) (ConversationSummary, error) {
	var summary ConversationSummary

	var targetID, targetPhoto string
	err := db.c.QueryRow("SELECT id, photo_url FROM users WHERE username = ?", targetUsername).
		Scan(&targetID, &targetPhoto)
	if errors.Is(err, sql.ErrNoRows) {
		return summary, ErrUserNotFound
	} else if err != nil {
		return summary, fmt.Errorf("error finding target user: %w", err)
	}

	if currentUserID == targetID {
		return summary, errors.New("cannot create conversation with yourself")
	}

	var existingID string
	err = db.c.QueryRow(`
		SELECT c.id FROM conversations c
		JOIN conversation_members cm1 ON c.id = cm1.conversation_id AND cm1.user_id = ?
		JOIN conversation_members cm2 ON c.id = cm2.conversation_id AND cm2.user_id = ?
		WHERE c.is_group = 0`, currentUserID, targetID).Scan(&existingID)
	if err == nil {
		summary.ID = existingID
		summary.Username = targetUsername
		summary.Photo = targetPhoto
		return summary, nil
	} else if !errors.Is(err, sql.ErrNoRows) {
		return summary, fmt.Errorf("error checking existing conversation: %w", err)
	}

	cUUID, err := uuid.NewV4()
	if err != nil {
		return summary, fmt.Errorf("error generating conversation UUID: %w", err)
	}
	convID := "conv_" + cUUID.String()

	tx, err := db.c.Begin()
	if err != nil {
		return summary, err
	}
	defer func() { _ = tx.Rollback() }()

	if _, err = tx.Exec("INSERT INTO conversations (id, is_group) VALUES (?, 0)", convID); err != nil {
		return summary, fmt.Errorf("error creating conversation: %w", err)
	}
	if _, err = tx.Exec(
		"INSERT INTO conversation_members (conversation_id, user_id) VALUES (?, ?), (?, ?)",
		convID, currentUserID, convID, targetID,
	); err != nil {
		return summary, fmt.Errorf("error adding conversation members: %w", err)
	}
	if err = tx.Commit(); err != nil {
		return summary, err
	}

	summary.ID = convID
	summary.Username = targetUsername
	summary.Photo = targetPhoto
	return summary, nil
}

func (db *appdbimpl) GetUserConversations(userID string) ([]ConversationSummary, error) {
	if err := db.markDelivered(userID); err != nil {
		return nil, err
	}

	rows, err := db.c.Query(`
		SELECT
			c.id,
			c.is_group,
			CASE WHEN c.is_group THEN c.name ELSE COALESCE((
				SELECT u.username FROM conversation_members cm2
				JOIN users u ON u.id = cm2.user_id
				WHERE cm2.conversation_id = c.id AND cm2.user_id != ?
				LIMIT 1
			), '') END AS title,
			CASE WHEN c.is_group THEN c.photo_url ELSE COALESCE((
				SELECT u.photo_url FROM conversation_members cm2
				JOIN users u ON u.id = cm2.user_id
				WHERE cm2.conversation_id = c.id AND cm2.user_id != ?
				LIMIT 1
			), '') END AS photo
		FROM conversations c
		JOIN conversation_members cm ON cm.conversation_id = c.id AND cm.user_id = ?`,
		userID, userID, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("error querying conversations: %w", err)
	}
	defer func() { _ = rows.Close() }()

	convs := []ConversationSummary{}
	for rows.Next() {
		var cs ConversationSummary
		if err := rows.Scan(&cs.ID, &cs.IsGroup, &cs.Username, &cs.Photo); err != nil {
			return nil, err
		}
		convs = append(convs, cs)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating conversations: %w", err)
	}

	for i := range convs {
		lm, err := db.lastMessage(convs[i].ID)
		if err != nil {
			return nil, err
		}
		convs[i].LastMessage = lm
	}

	sort.SliceStable(convs, func(i, j int) bool {
		ti, tj := time.Time{}, time.Time{}
		if convs[i].LastMessage != nil {
			ti = convs[i].LastMessage.DataSent
		}
		if convs[j].LastMessage != nil {
			tj = convs[j].LastMessage.DataSent
		}
		return ti.After(tj)
	})

	return convs, nil
}

func (db *appdbimpl) lastMessage(convID string) (*LastMessage, error) {
	var (
		text, photo, sender string
		dataSent            time.Time
	)
	err := db.c.QueryRow(`
		SELECT m.content_text, m.content_photo, m.data_sent, COALESCE(u.username, '')
		FROM messages m
		LEFT JOIN users u ON u.id = m.sender_id
		WHERE m.conversation_id = ?
		ORDER BY m.data_sent DESC, m.id DESC
		LIMIT 1`, convID).Scan(&text, &photo, &dataSent, &sender)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error reading last message: %w", err)
	}

	return &LastMessage{
		Text:           text,
		IsPhoto:        photo != "",
		DataSent:       dataSent,
		SenderUsername: sender,
	}, nil
}

func (db *appdbimpl) markDelivered(userID string) error {
	_, err := db.c.Exec(`
		INSERT OR IGNORE INTO message_status (message_id, user_id, status)
		SELECT m.id, ?, 'delivered'
		FROM messages m
		JOIN conversation_members cm ON cm.conversation_id = m.conversation_id AND cm.user_id = ?
		WHERE m.sender_id != ?`, userID, userID, userID)
	if err != nil {
		return fmt.Errorf("error marking messages delivered: %w", err)
	}
	return nil
}

func (db *appdbimpl) markRead(convID, userID string) error {
	_, err := db.c.Exec(`
		INSERT INTO message_status (message_id, user_id, status)
		SELECT m.id, ?, 'read'
		FROM messages m
		WHERE m.conversation_id = ? AND m.sender_id != ?
		ON CONFLICT(message_id, user_id) DO UPDATE SET status = 'read'`,
		userID, convID, userID)
	if err != nil {
		return fmt.Errorf("error marking messages read: %w", err)
	}
	return nil
}

func (db *appdbimpl) messageStatus(convID, msgID, senderID string) (string, error) {
	var total, readCnt, delivCnt int
	err := db.c.QueryRow(`
		SELECT
			COUNT(*),
			COALESCE(SUM(CASE WHEN ms.status = 'read' THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN ms.status IN ('read', 'delivered') THEN 1 ELSE 0 END), 0)
		FROM conversation_members cm
		LEFT JOIN message_status ms ON ms.message_id = ? AND ms.user_id = cm.user_id
		WHERE cm.conversation_id = ? AND cm.user_id != ?`,
		msgID, convID, senderID).Scan(&total, &readCnt, &delivCnt)
	if err != nil {
		return "", fmt.Errorf("error computing message status: %w", err)
	}

	switch {
	case total > 0 && readCnt == total:
		return "read", nil
	case total > 0 && delivCnt == total:
		return "delivered", nil
	default:
		return "sent", nil
	}
}

func (db *appdbimpl) GetConversationDetails(convID, userID string) (ConversationDetails, error) {
	member, err := db.isMember(convID, userID)
	if err != nil {
		return ConversationDetails{}, err
	}
	if !member {
		return ConversationDetails{}, ErrNotFound
	}

	if err := db.markRead(convID, userID); err != nil {
		return ConversationDetails{}, err
	}

	var details ConversationDetails
	details.ID = convID
	details.Username, details.Photo, details.IsGroup, err = db.conversationView(convID, userID)
	if err != nil {
		return ConversationDetails{}, fmt.Errorf("error reading conversation: %w", err)
	}

	rows, err := db.c.Query(`
		SELECT m.id, m.conversation_id, m.sender_id, COALESCE(u.username, ''),
		       COALESCE(u.photo_url, ''),
		       m.content_text, m.content_photo, COALESCE(m.reply_to_message_id, ''),
		       m.is_forwarded, m.data_sent
		FROM messages m
		LEFT JOIN users u ON u.id = m.sender_id
		WHERE m.conversation_id = ?
		ORDER BY m.data_sent ASC, m.id ASC`, convID)
	if err != nil {
		return ConversationDetails{}, fmt.Errorf("error querying messages: %w", err)
	}
	defer func() { _ = rows.Close() }()

	details.Messages = []Message{}
	for rows.Next() {
		var m Message
		if err := rows.Scan(&m.ID, &m.ConversationID, &m.SenderID, &m.SenderUsername,
			&m.SenderPhotoURL, &m.ContentText, &m.ContentPhoto, &m.ReplyToMessageID, &m.IsForwarded, &m.DataSent); err != nil {
			return ConversationDetails{}, err
		}
		details.Messages = append(details.Messages, m)
	}
	if err := rows.Err(); err != nil {
		return ConversationDetails{}, fmt.Errorf("error iterating messages: %w", err)
	}

	for i := range details.Messages {
		m := &details.Messages[i]
		if m.Status, err = db.messageStatus(convID, m.ID, m.SenderID); err != nil {
			return ConversationDetails{}, err
		}
		if m.Reactions, err = db.messageReactions(m.ID); err != nil {
			return ConversationDetails{}, err
		}
	}

	return details, nil
}
