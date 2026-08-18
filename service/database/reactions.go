package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/gofrs/uuid"
)

type Reaction struct {
	ReactionType string `json:"reactionType"`
	UserSenderID string `json:"userSenderId"`
}

// AddReaction aggiunge una reazione a un messaggio
func (db *appdbimpl) AddReaction(messageID, userID, reactionType string) (Reaction, error) {
	var r Reaction

	// Verifica che il messaggio esista
	var exists bool
	err := db.c.QueryRow("SELECT EXISTS(SELECT 1 FROM messages WHERE id = ?)", messageID).Scan(&exists)
	if err != nil || !exists {
		return r, errors.New("message not found")
	}

	rUUID, err := uuid.NewV4()
	if err != nil {
		return r, fmt.Errorf("error generating reaction UUID: %w", err)
	}
	reactionID := "react_" + rUUID.String()

	query := "INSERT INTO reactions (id, message_id, user_id, reaction_type) VALUES (?, ?, ?, ?)"
	_, err = db.c.Exec(query, reactionID, messageID, userID, reactionType)
	if err != nil {
		return r, fmt.Errorf("error inserting reaction: %w", err)
	}

	r.ReactionType = reactionType
	r.UserSenderID = userID

	return r, nil
}

// RemoveReaction rimuove una reazione da un messaggio
func (db *appdbimpl) RemoveReaction(messageID, reactionID, userID string) error {
	res, err := db.c.Exec("DELETE FROM reactions WHERE id = ? AND message_id = ? AND user_id = ?", reactionID, messageID, userID)
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("reaction not found or access denied")
	}

	return nil
}

// ForwardMessage copia un messaggio esistente in una nuova conversazione
func (db *appdbimpl) ForwardMessage(originalMsgID, targetConvID, senderUsername string) (Message, error) {
	var msg Message

	// Recupera il contenuto del messaggio originale
	var cType, cVal string
	err := db.c.QueryRow("SELECT content_type, content_value FROM messages WHERE id = ?", originalMsgID).Scan(&cType, &cVal)
	if errors.Is(err, sql.ErrNoRows) {
		return msg, errors.New("original message not found")
	} else if err != nil {
		return msg, err
	}

	// Invia il messaggio recuperato nella conversazione di destinazione
	return db.SendMessage(targetConvID, senderUsername, cType, cVal)
}
