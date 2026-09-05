package database

import (
	"fmt"

	"github.com/gofrs/uuid"
)

func (db *appdbimpl) messageReactions(messageID string) ([]Reaction, error) {
	rows, err := db.c.Query(`
		SELECT r.id, r.reaction_type, r.user_id, COALESCE(u.username, '')
		FROM reactions r
		LEFT JOIN users u ON u.id = r.user_id
		WHERE r.message_id = ?
		ORDER BY r.id`, messageID)
	if err != nil {
		return nil, fmt.Errorf("error querying reactions: %w", err)
	}
	defer func() { _ = rows.Close() }()

	reactions := []Reaction{}
	for rows.Next() {
		var r Reaction
		if err := rows.Scan(&r.ID, &r.ReactionType, &r.UserSenderID, &r.Username); err != nil {
			return nil, err
		}
		reactions = append(reactions, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return reactions, nil
}

func (db *appdbimpl) messageInConversation(convID, messageID, userID string) (bool, error) {
	var ok bool
	err := db.c.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM messages m
			JOIN conversation_members cm ON cm.conversation_id = m.conversation_id AND cm.user_id = ?
			WHERE m.id = ? AND m.conversation_id = ?
		)`, userID, messageID, convID).Scan(&ok)
	return ok, err
}

func (db *appdbimpl) AddReaction(convID, messageID, userID, reactionType string) (Reaction, error) {
	var r Reaction

	ok, err := db.messageInConversation(convID, messageID, userID)
	if err != nil {
		return r, err
	}
	if !ok {
		return r, ErrNotFound
	}

	rUUID, err := uuid.NewV4()
	if err != nil {
		return r, fmt.Errorf("error generating reaction UUID: %w", err)
	}
	reactionID := "react_" + rUUID.String()

	_, err = db.c.Exec(`
		INSERT INTO reactions (id, message_id, user_id, reaction_type)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(message_id, user_id) DO UPDATE SET reaction_type = excluded.reaction_type`,
		reactionID, messageID, userID, reactionType)
	if err != nil {
		return r, fmt.Errorf("error inserting reaction: %w", err)
	}

	var username string
	if err := db.c.QueryRow(`
		SELECT r.id, COALESCE(u.username, '')
		FROM reactions r LEFT JOIN users u ON u.id = r.user_id
		WHERE r.message_id = ? AND r.user_id = ?`, messageID, userID).Scan(&reactionID, &username); err != nil {
		return r, fmt.Errorf("error reading reaction: %w", err)
	}

	return Reaction{ID: reactionID, ReactionType: reactionType, UserSenderID: userID, Username: username}, nil
}

func (db *appdbimpl) RemoveReaction(convID, messageID, reactionID, userID string) error {
	ok, err := db.messageInConversation(convID, messageID, userID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrNotFound
	}

	res, err := db.c.Exec(
		"DELETE FROM reactions WHERE id = ? AND message_id = ? AND user_id = ?",
		reactionID, messageID, userID,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
