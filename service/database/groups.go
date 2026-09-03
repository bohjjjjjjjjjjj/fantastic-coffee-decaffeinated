package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/gofrs/uuid"
)

// CreateGroup crea un nuovo gruppo con l'utente creatore come membro. memberIDs
// (opzionale) elenca altri utenti da aggiungere subito; gli ID inesistenti sono
// ignorati.
func (db *appdbimpl) CreateGroup(groupName, ownerID string, memberIDs []string) (GroupResponse, error) {
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
	defer func() { _ = tx.Rollback() }()

	if _, err = tx.Exec("INSERT INTO conversations (id, name, is_group) VALUES (?, ?, 1)", groupID, groupName); err != nil {
		return g, fmt.Errorf("error inserting group conversation: %w", err)
	}
	if _, err = tx.Exec("INSERT INTO conversation_members (conversation_id, user_id) VALUES (?, ?)", groupID, ownerID); err != nil {
		return g, fmt.Errorf("error adding group owner: %w", err)
	}
	for _, memberID := range memberIDs {
		if memberID == "" || memberID == ownerID {
			continue
		}
		if _, err = tx.Exec(`
			INSERT OR IGNORE INTO conversation_members (conversation_id, user_id)
			SELECT ?, id FROM users WHERE id = ?`, groupID, memberID); err != nil {
			return g, fmt.Errorf("error adding group member: %w", err)
		}
	}
	if err = tx.Commit(); err != nil {
		return g, err
	}

	g.ID = groupID
	g.Name = groupName
	g.PhotoURL = ""
	return g, nil
}

// GetGroupDetails restituisce i dettagli di un gruppo di cui l'utente fa parte.
func (db *appdbimpl) GetGroupDetails(groupID, userID string) (GroupResponse, error) {
	var g GroupResponse

	member, err := db.isMember(groupID, userID)
	if err != nil {
		return g, err
	}
	if !member {
		return g, ErrNotFound
	}

	err = db.c.QueryRow(
		"SELECT id, name, photo_url FROM conversations WHERE id = ? AND is_group = 1",
		groupID,
	).Scan(&g.ID, &g.Name, &g.PhotoURL)
	if errors.Is(err, sql.ErrNoRows) {
		return g, ErrNotFound
	}
	if err != nil {
		return g, fmt.Errorf("error reading group: %w", err)
	}
	return g, nil
}

// GetGroupMembers restituisce l'elenco dei membri di un gruppo di cui l'utente fa
// parte.
func (db *appdbimpl) GetGroupMembers(groupID, userID string) ([]User, error) {
	member, err := db.isMember(groupID, userID)
	if err != nil {
		return nil, err
	}
	if !member {
		return nil, ErrNotFound
	}
	return db.groupMembers(groupID)
}

func (db *appdbimpl) groupMembers(groupID string) ([]User, error) {
	rows, err := db.c.Query(`
		SELECT u.id, u.username, u.photo_url
		FROM conversation_members cm
		JOIN users u ON u.id = cm.user_id
		WHERE cm.conversation_id = ?
		ORDER BY u.username`, groupID)
	if err != nil {
		return nil, fmt.Errorf("error querying group members: %w", err)
	}
	defer func() { _ = rows.Close() }()

	members := []User{}
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Username, &u.PhotoURL); err != nil {
			return nil, err
		}
		members = append(members, u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return members, nil
}

// AddToGroup aggiunge membri a un gruppo di cui l'utente fa parte e restituisce
// l'elenco aggiornato. Gli ID inesistenti sono ignorati.
func (db *appdbimpl) AddToGroup(groupID, userID string, memberIDs []string) ([]User, error) {
	member, err := db.isMember(groupID, userID)
	if err != nil {
		return nil, err
	}
	if !member {
		return nil, ErrNotFound
	}

	var isGroup bool
	if err := db.c.QueryRow("SELECT is_group FROM conversations WHERE id = ?", groupID).Scan(&isGroup); err != nil || !isGroup {
		return nil, ErrNotFound
	}

	for _, memberID := range memberIDs {
		if memberID == "" {
			continue
		}
		if _, err := db.c.Exec(`
			INSERT OR IGNORE INTO conversation_members (conversation_id, user_id)
			SELECT ?, id FROM users WHERE id = ?`, groupID, memberID); err != nil {
			return nil, fmt.Errorf("error adding group member: %w", err)
		}
	}

	return db.groupMembers(groupID)
}

// SetGroupName modifica il nome di un gruppo di cui l'utente fa parte.
func (db *appdbimpl) SetGroupName(groupID, newName, userID string) error {
	member, err := db.isMember(groupID, userID)
	if err != nil {
		return err
	}
	if !member {
		return ErrNotFound
	}

	res, err := db.c.Exec("UPDATE conversations SET name = ? WHERE id = ? AND is_group = 1", newName, groupID)
	if err != nil {
		return fmt.Errorf("error updating group name: %w", err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrNotFound
	}
	return nil
}

// SetGroupPhoto modifica la foto di un gruppo di cui l'utente fa parte.
func (db *appdbimpl) SetGroupPhoto(groupID, photoURL, userID string) error {
	member, err := db.isMember(groupID, userID)
	if err != nil {
		return err
	}
	if !member {
		return ErrNotFound
	}

	res, err := db.c.Exec("UPDATE conversations SET photo_url = ? WHERE id = ? AND is_group = 1", photoURL, groupID)
	if err != nil {
		return fmt.Errorf("error updating group photo: %w", err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrNotFound
	}
	return nil
}

// LeaveGroup rimuove l'utente dai membri del gruppo. Se il gruppo resta senza
// membri viene eliminato.
func (db *appdbimpl) LeaveGroup(groupID, userID string) error {
	tx, err := db.c.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	res, err := tx.Exec(
		"DELETE FROM conversation_members WHERE conversation_id = ? AND user_id = ? AND conversation_id IN (SELECT id FROM conversations WHERE is_group = 1)",
		groupID, userID,
	)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrNotFound
	}

	var remaining int
	if err := tx.QueryRow("SELECT COUNT(*) FROM conversation_members WHERE conversation_id = ?", groupID).Scan(&remaining); err != nil {
		return err
	}
	if remaining == 0 {
		// Rimozione esplicita delle righe collegate: non si fa affidamento sulle
		// FOREIGN KEY cascade, che con il pool di connessioni non sono garantite.
		stmts := []string{
			"DELETE FROM reactions WHERE message_id IN (SELECT id FROM messages WHERE conversation_id = ?)",
			"DELETE FROM message_status WHERE message_id IN (SELECT id FROM messages WHERE conversation_id = ?)",
			"DELETE FROM messages WHERE conversation_id = ?",
			"DELETE FROM conversation_members WHERE conversation_id = ?",
			"DELETE FROM conversations WHERE id = ?",
		}
		for _, s := range stmts {
			if _, err := tx.Exec(s, groupID); err != nil {
				return fmt.Errorf("error deleting empty group: %w", err)
			}
		}
	}

	return tx.Commit()
}
