package database

import (
	"database/sql"
	"errors"
	"fmt"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/globaltime"
	"github.com/gofrs/uuid"
)

// SendMessage salva un nuovo messaggio nella conversazione. L'utente deve farne
// parte, altrimenti viene restituito ErrNotFound.
func (db *appdbimpl) SendMessage(convID, senderID, cType, cVal, replyToID string) (Message, error) {
	var msg Message

	member, err := db.isMember(convID, senderID)
	if err != nil {
		return msg, err
	}
	if !member {
		return msg, ErrNotFound
	}

	if replyToID != "" {
		var ok bool
		if err := db.c.QueryRow(
			`SELECT EXISTS(SELECT 1 FROM messages WHERE id = ? AND conversation_id = ?)`,
			replyToID, convID,
		).Scan(&ok); err != nil {
			return msg, fmt.Errorf("error checking reply target: %w", err)
		}
		if !ok {
			return msg, ErrNotFound
		}
	}

	mUUID, err := uuid.NewV4()
	if err != nil {
		return msg, fmt.Errorf("error generating message UUID: %w", err)
	}
	msgID := "msg_" + mUUID.String()
	now := globaltime.Now().UTC()

	var replyValue interface{}
	if replyToID != "" {
		replyValue = replyToID
	}

	_, err = db.c.Exec(`
		INSERT INTO messages (id, conversation_id, sender_id, content_type, content_value, reply_to_message_id, data_sent)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		msgID, convID, senderID, cType, cVal, replyValue, now)
	if err != nil {
		return msg, fmt.Errorf("error inserting message: %w", err)
	}

	var senderUsername string
	if err := db.c.QueryRow("SELECT username FROM users WHERE id = ?", senderID).Scan(&senderUsername); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return msg, fmt.Errorf("error reading sender: %w", err)
	}

	return Message{
		ID:               msgID,
		ConversationID:   convID,
		SenderID:         senderID,
		SenderUsername:   senderUsername,
		ContentType:      cType,
		ContentValue:     cVal,
		ReplyToMessageID: replyToID,
		DataSent:         now,
		Status:           "sent",
		Reactions:        []Reaction{},
	}, nil
}

// ForwardMessage copia un messaggio in un'altra conversazione. L'utente deve far
// parte sia della conversazione di origine sia di quella di destinazione.
func (db *appdbimpl) ForwardMessage(originalMsgID, targetConvID, senderID string) (Message, error) {
	var msg Message

	var srcConvID, cType, cVal string
	err := db.c.QueryRow(
		"SELECT conversation_id, content_type, content_value FROM messages WHERE id = ?",
		originalMsgID,
	).Scan(&srcConvID, &cType, &cVal)
	if errors.Is(err, sql.ErrNoRows) {
		return msg, ErrNotFound
	} else if err != nil {
		return msg, fmt.Errorf("error reading original message: %w", err)
	}

	member, err := db.isMember(srcConvID, senderID)
	if err != nil {
		return msg, err
	}
	if !member {
		return msg, ErrNotFound
	}

	return db.SendMessage(targetConvID, senderID, cType, cVal, "")
}

// DeleteMessage elimina un messaggio inviato dall'utente. La proprietà è legata
// a sender_id, quindi resta valida anche dopo un cambio di username.
func (db *appdbimpl) DeleteMessage(msgID, convID, senderID string) error {
	tx, err := db.c.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	res, err := tx.Exec(
		"DELETE FROM messages WHERE id = ? AND conversation_id = ? AND sender_id = ?",
		msgID, convID, senderID,
	)
	if err != nil {
		return err
	}
	if n, err := res.RowsAffected(); err != nil || n == 0 {
		return ErrNotFound
	}

	if _, err := tx.Exec("DELETE FROM reactions WHERE message_id = ?", msgID); err != nil {
		return err
	}
	if _, err := tx.Exec("DELETE FROM message_status WHERE message_id = ?", msgID); err != nil {
		return err
	}

	return tx.Commit()
}
