package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/gofrs/uuid"
)

// ConversationSummary rappresenta la sintesi della chat per la lista
type ConversationSummary struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Photo    string `json:"photo"`
}

// GetOrCreateConversation crea o recupera una chat tra due utenti
func (db *appdbimpl) CreateConversation(currentUserID, targetUsername string) (ConversationSummary, error) {
	var summary ConversationSummary

	// 1. Cerca l'utente destinatario tramite username
	var targetID, targetPhoto string
	err := db.c.QueryRow("SELECT id, photo_url FROM users WHERE username = ?", targetUsername).Scan(&targetID, &targetPhoto)
	if errors.Is(err, sql.ErrNoRows) {
		return summary, ErrUserNotFound
	} else if err != nil {
		return summary, fmt.Errorf("error finding target user: %w", err)
	}

	// Non si può creare una chat con se stessi
	if currentUserID == targetID {
		return summary, errors.New("cannot create conversation with yourself")
	}

	// 2. Controlla se esiste già una conversazione diretta (non gruppo) tra questi due utenti
	queryCheck := `
		SELECT c.id FROM conversations c
		JOIN conversation_members cm1 ON c.id = cm1.conversation_id
		JOIN conversation_members cm2 ON c.id = cm2.conversation_id
		WHERE c.is_group = FALSE AND cm1.user_id = ? AND cm2.user_id = ?`

	var existingID string
	err = db.c.QueryRow(queryCheck, currentUserID, targetID).Scan(&existingID)
	if err == nil {
		// Conversazione già esistente
		summary.ID = existingID
		summary.Username = targetUsername
		summary.Photo = targetPhoto
		return summary, nil
	}

	// 3. Se non esiste, crea un nuovo UUID per la conversazione
	cUUID, err := uuid.NewV4()
	if err != nil {
		return summary, fmt.Errorf("error generating conversation UUID: %w", err)
	}
	convID := "conv_" + cUUID.String()

	// Inserimento transazionale nel DB
	tx, err := db.c.Begin()
	if err != nil {
		return summary, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	_, err = tx.Exec("INSERT INTO conversations (id, is_group) VALUES (?, FALSE)", convID)
	if err != nil {
		return summary, fmt.Errorf("error creating conversation: %w", err)
	}

	// Aggiungi entrambi gli utenti come membri
	_, err = tx.Exec("INSERT INTO conversation_members (conversation_id, user_id) VALUES (?, ?), (?, ?)",
		convID, currentUserID, convID, targetID)
	if err != nil {
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

// GetUserConversations recupera tutte le chat dell'utente autenticato
func (db *appdbimpl) GetUserConversations(userID string) ([]ConversationSummary, error) {
	query := `
		SELECT c.id, 
		       COALESCE(u.username, c.name) AS name, 
		       COALESCE(u.photo_url, c.photo_url) AS photo
		FROM conversations c
		JOIN conversation_members cm ON c.id = cm.conversation_id
		LEFT JOIN conversation_members cm_other ON c.id = cm_other.conversation_id AND cm_other.user_id != ?
		LEFT JOIN users u ON cm_other.user_id = u.id
		WHERE cm.user_id = ?`

	rows, err := db.c.Query(query, userID, userID)
	if err != nil {
		return nil, fmt.Errorf("error querying conversations: %w", err)
	}
	defer rows.Close()

	var convs []ConversationSummary
	for rows.Next() {
		var cs ConversationSummary
		if err := rows.Scan(&cs.ID, &cs.Username, &cs.Photo); err != nil {
			return nil, err
		}
		convs = append(convs, cs)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating conversations: %w", err)
	}

	if convs == nil {
		convs = []ConversationSummary{}
	}

	return convs, nil
}
