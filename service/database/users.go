package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/gofrs/uuid"
)

// GetOrCreateUser cerca l'utente o lo crea se non esiste, generando un nuovo token di sessione.
func (db *appdbimpl) GetOrCreateUser(username string) (string, string, error) {
	var userID string

	// 1. Cerca se l'utente esiste già
	err := db.c.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&userID)
	if errors.Is(err, sql.ErrNoRows) {
		// Se non esiste, crea un nuovo UUID per l'utente
		uUUID, err := uuid.NewV4()
		if err != nil {
			return "", "", fmt.Errorf("error generating user UUID: %w", err)
		}
		userID = "usr_" + uUUID.String()

		_, err = db.c.Exec("INSERT INTO users (id, username) VALUES (?, ?)", userID, username)
		if err != nil {
			return "", "", fmt.Errorf("error creating new user: %w", err)
		}
	} else if err != nil {
		return "", "", fmt.Errorf("database query error: %w", err)
	}

	// 2. Genera un token di sessione univoco
	sUUID, err := uuid.NewV4()
	if err != nil {
		return "", "", fmt.Errorf("error generating session UUID: %w", err)
	}
	sessionToken := "sess_" + sUUID.String()

	// 3. Salva la sessione nel DB
	_, err = db.c.Exec("INSERT INTO sessions (token, user_id) VALUES (?, ?)", sessionToken, userID)
	if err != nil {
		return "", "", fmt.Errorf("error saving session: %w", err)
	}

	return userID, sessionToken, nil
}

// GetUserByToken recupera le info dell'utente tramite il token Bearer dell'header Authorization
func (db *appdbimpl) GetUserByToken(token string) (string, string, error) {
	var userID, username string

	query := `
		SELECT u.id, u.username 
		FROM users u 
		JOIN sessions s ON u.id = s.user_id 
		WHERE s.token = ?`

	err := db.c.QueryRow(query, token).Scan(&userID, &username)
	if errors.Is(err, sql.ErrNoRows) {
		return "", "", ErrUserNotFound
	} else if err != nil {
		return "", "", fmt.Errorf("error validating session token: %w", err)
	}

	return userID, username, nil
}

// GetUserByID restituisce i dettagli di un utente tramite il suo ID
func (db *appdbimpl) GetUserByID(userID string) (User, error) {
	var u User
	query := "SELECT id, username, photo_url FROM users WHERE id = ?"
	err := db.c.QueryRow(query, userID).Scan(&u.ID, &u.Username, &u.PhotoURL)
	if errors.Is(err, sql.ErrNoRows) {
		return u, ErrUserNotFound
	}
	if err != nil {
		return u, fmt.Errorf("error getting user by id: %w", err)
	}
	return u, nil
}

// UpdateUsername aggiorna lo username dell'utente nel DB
func (db *appdbimpl) UpdateUsername(userID string, newUsername string) error {
	_, err := db.c.Exec("UPDATE users SET username = ? WHERE id = ?", newUsername, userID)
	if err != nil {
		return fmt.Errorf("error updating username: %w", err)
	}
	return nil
}

// SearchUsers cerca gli utenti il cui username contiene la stringa fornita
func (db *appdbimpl) SearchUsers(searchQuery string) ([]User, error) {
	query := "SELECT id, username, photo_url FROM users WHERE username LIKE ? LIMIT 100"
	rows, err := db.c.Query(query, "%"+searchQuery+"%")
	if err != nil {
		return nil, fmt.Errorf("error searching users: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Username, &u.PhotoURL); err != nil {
			return nil, fmt.Errorf("error scanning user row: %w", err)
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	// Se non trova nessuno, restituisce una slice vuota anziché nil (per produrre "[]" in JSON invece di "null")
	if users == nil {
		users = []User{}
	}

	return users, nil
}
