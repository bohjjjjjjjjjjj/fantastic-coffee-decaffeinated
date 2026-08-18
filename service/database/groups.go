package database

import (
	// "database/sql"
	"errors"
	"fmt"

	"github.com/gofrs/uuid"
)

type GroupResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	PhotoURL string `json:"photoUrl"`
}

// CreateGroup crea una nuova chat di gruppo con l'utente creatore come membro
func (db *appdbimpl) CreateGroup(groupName, ownerID string) (GroupResponse, error) {
	var g GroupResponse

	gUUID, err := uuid.NewV4()
	if err != nil {
		return g, fmt.Errorf("error generating group UUID: %w", err)
	}
	groupID := "grp_" + gUUID.String()

	tx, err := db.c.Begin()
	if err != nil {
		return g, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	// Inserisce la conversazione segnata come gruppo (is_group = TRUE)
	_, err = tx.Exec("INSERT INTO conversations (id, name, is_group) VALUES (?, ?, TRUE)", groupID, groupName)
	if err != nil {
		return g, fmt.Errorf("error inserting group conversation: %w", err)
	}

	// Aggiunge il creatore ai membri
	_, err = tx.Exec("INSERT INTO conversation_members (conversation_id, user_id) VALUES (?, ?)", groupID, ownerID)
	if err != nil {
		return g, fmt.Errorf("error adding group owner: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return g, err
	}

	g.ID = groupID
	g.Name = groupName
	g.PhotoURL = ""

	return g, nil
}

// SetGroupName modifica il nome del gruppo
func (db *appdbimpl) SetGroupName(groupID, newName, userID string) error {
	// Verifica che il gruppo esista e che l'utente sia un membro
	var isMember bool
	err := db.c.QueryRow("SELECT EXISTS(SELECT 1 FROM conversation_members WHERE conversation_id = ? AND user_id = ?)", groupID, userID).Scan(&isMember)
	if err != nil || !isMember {
		return errors.New("group not found or user is not a member")
	}

	_, err = db.c.Exec("UPDATE conversations SET name = ? WHERE id = ? AND is_group = TRUE", newName, groupID)
	if err != nil {
		return fmt.Errorf("error updating group name: %w", err)
	}

	return nil
}

// LeaveGroup rimuove l'utente corrente dai membri del gruppo
func (db *appdbimpl) LeaveGroup(groupID, userID string) error {
	res, err := db.c.Exec("DELETE FROM conversation_members WHERE conversation_id = ? AND user_id = ?", groupID, userID)
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("membership not found")
	}

	return nil
}
