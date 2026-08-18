package database

import (
	"errors"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
)

type Message struct {
	ID             string    `json:"id"`
	ConversationID string    `json:"conversationId"`
	SenderUsername string    `json:"senderUsername"`
	ContentType    string    `json:"contentType"`
	ContentValue   string    `json:"contentValue"`
	Timestamp      time.Time `json:"timestamp"`
}

type ConversationDetails struct {
	ID       string    `json:"id"`
	Messages []Message `json:"messages"`
}

// SendMessage salva un nuovo messaggio nel database
func (db *appdbimpl) SendMessage(convID, senderUsername, cType, cVal string) (Message, error) {
	var msg Message

	mUUID, err := uuid.NewV4()
	if err != nil {
		return msg, fmt.Errorf("error generating message UUID: %w", err)
	}

	msgID := "msg_" + mUUID.String()
	now := time.Now()

	query := `
		INSERT INTO messages (id, conversation_id, sender_username, content_type, content_value, timestamp)
		VALUES (?, ?, ?, ?, ?, ?)`

	_, err = db.c.Exec(query, msgID, convID, senderUsername, cType, cVal, now)
	if err != nil {
		return msg, fmt.Errorf("error inserting message: %w", err)
	}

	msg = Message{
		ID:             msgID,
		ConversationID: convID,
		SenderUsername: senderUsername,
		ContentType:    cType,
		ContentValue:   cVal,
		Timestamp:      now,
	}

	return msg, nil
}

// GetConversationDetails recupera la conversazione e tutti i suoi messaggi
func (db *appdbimpl) GetConversationDetails(convID, userID string) (ConversationDetails, error) {
	var details ConversationDetails

	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM conversation_members WHERE conversation_id = ? AND user_id = ?", convID, userID).Scan(&count)
	if err != nil || count == 0 {
		return ConversationDetails{}, errors.New("access denied or conversation not found")
	}

	details.ID = convID

	queryMsgs := `
		SELECT id, conversation_id, sender_username, content_type, content_value, timestamp
		FROM messages
		WHERE conversation_id = ?
		ORDER BY timestamp ASC`

	rows, err := db.c.Query(queryMsgs, convID)
	if err != nil {
		return ConversationDetails{}, fmt.Errorf("error querying messages: %w", err)
	}
	defer rows.Close()

	var msgs []Message
	for rows.Next() {
		var m Message
		if err := rows.Scan(&m.ID, &m.ConversationID, &m.SenderUsername, &m.ContentType, &m.ContentValue, &m.Timestamp); err != nil {
			return ConversationDetails{}, err
		}
		msgs = append(msgs, m)
	}

	if err := rows.Err(); err != nil {
		return ConversationDetails{}, fmt.Errorf("error iterating messages: %w", err)
	}

	if msgs == nil {
		msgs = []Message{}
	}

	details.Messages = msgs
	return details, nil
}

// DeleteMessage elimina un messaggio Inviato dall'utente
func (db *appdbimpl) DeleteMessage(msgID, convID, username string) error {
	res, err := db.c.Exec("DELETE FROM messages WHERE id = ? AND conversation_id = ? AND sender_username = ?", msgID, convID, username)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return errors.New("message not found or not authorized")
	}

	return nil
}
