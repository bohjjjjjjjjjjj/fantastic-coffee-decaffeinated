package database

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/mattn/go-sqlite3"
)

func (db *appdbimpl) GetOrCreateUser(username string) (string, string, bool, error) {
	var userID string
	created := false

	err := db.c.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&userID)
	if errors.Is(err, sql.ErrNoRows) {
		uUUID, err := uuid.NewV4()
		if err != nil {
			return "", "", false, fmt.Errorf("error generating user UUID: %w", err)
		}
		userID = "usr_" + uUUID.String()

		if _, err = db.c.Exec("INSERT INTO users (id, username) VALUES (?, ?)", userID, username); err != nil {
			return "", "", false, fmt.Errorf("error creating new user: %w", err)
		}
		created = true
	} else if err != nil {
		return "", "", false, fmt.Errorf("database query error: %w", err)
	}

	sUUID, err := uuid.NewV4()
	if err != nil {
		return "", "", false, fmt.Errorf("error generating session UUID: %w", err)
	}
	sessionToken := "sess_" + sUUID.String()

	if _, err = db.c.Exec("INSERT INTO sessions (token, user_id) VALUES (?, ?)", sessionToken, userID); err != nil {
		return "", "", false, fmt.Errorf("error saving session: %w", err)
	}

	return userID, sessionToken, created, nil
}

func (db *appdbimpl) GetUserByToken(token string) (string, string, error) {
	var userID, username string
	err := db.c.QueryRow(`
		SELECT u.id, u.username
		FROM users u
		JOIN sessions s ON u.id = s.user_id
		WHERE s.token = ?`, token).Scan(&userID, &username)
	if errors.Is(err, sql.ErrNoRows) {
		return "", "", ErrUserNotFound
	} else if err != nil {
		return "", "", fmt.Errorf("error validating session token: %w", err)
	}
	return userID, username, nil
}

func (db *appdbimpl) GetUserByID(userID string) (User, error) {
	var u User
	err := db.c.QueryRow("SELECT id, username, photo_url FROM users WHERE id = ?", userID).
		Scan(&u.ID, &u.Username, &u.PhotoURL)
	if errors.Is(err, sql.ErrNoRows) {
		return u, ErrUserNotFound
	}
	if err != nil {
		return u, fmt.Errorf("error getting user by id: %w", err)
	}
	return u, nil
}

func (db *appdbimpl) UpdateUsername(userID string, newUsername string) error {
	_, err := db.c.Exec("UPDATE users SET username = ? WHERE id = ?", newUsername, userID)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return ErrUsernameTaken
		}
		return fmt.Errorf("error updating username: %w", err)
	}
	return nil
}

func (db *appdbimpl) SetUserPhoto(userID string, photoURL string) error {
	if _, err := db.c.Exec("UPDATE users SET photo_url = ? WHERE id = ?", photoURL, userID); err != nil {
		return fmt.Errorf("error updating user photo: %w", err)
	}
	return nil
}

func (db *appdbimpl) SearchUsers(searchQuery string, excludeUserID string) ([]User, error) {
	esc := strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(searchQuery)

	rows, err := db.c.Query(
		`SELECT id, username, photo_url FROM users
		 WHERE id != ? AND username LIKE ? ESCAPE '\'
		 ORDER BY username LIMIT 100`,
		excludeUserID, "%"+esc+"%",
	)
	if err != nil {
		return nil, fmt.Errorf("error searching users: %w", err)
	}
	defer func() { _ = rows.Close() }()

	users := []User{}
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
	return users, nil
}
