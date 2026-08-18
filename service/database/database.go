package database

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// ErrUserNotFound viene restituito quando un utente non esiste
var ErrUserNotFound = errors.New("user not found")

// User rappresenta un utente nel database
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	PhotoURL string `json:"photoUrl"`
}

// AppDatabase è l'interfaccia che definisce le operazioni sul database
type AppDatabase interface {
	// Autenticazione e Sessioni
	GetOrCreateUser(username string) (string, string, error) // Ritorna (userID, sessionToken, error)
	GetUserByToken(token string) (string, string, error)     // Ritorna (userID, username, error)

	// Profilo Utente e Ricerca (PASSO 1)
	GetUserByID(userID string) (User, error)
	UpdateUsername(userID string, newUsername string) error
	SearchUsers(searchQuery string) ([]User, error)

	CreateConversation(currentUserID, targetUsername string) (ConversationSummary, error)
	GetUserConversations(userID string) ([]ConversationSummary, error)

	SendMessage(convID, senderUsername, contentType, contentValue string) (Message, error)
	GetConversationDetails(convID, userID string) (ConversationDetails, error)
	DeleteMessage(msgID, convID, senderUsername string) error

	CreateGroup(groupName, ownerID string) (GroupResponse, error)
	SetGroupName(groupID, newName, userID string) error
	LeaveGroup(groupID, userID string) error

	AddReaction(messageID, userID, reactionType string) (Reaction, error)
	RemoveReaction(messageID, reactionID, userID string) error
	ForwardMessage(originalMsgID, targetConvID, senderUsername string) (Message, error)
	// Controllo stato
	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New crea e inizializza una nuova connessione al database SQLite
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("a database is required")
	}

	// Creazione delle tabelle se non esistono
	var tableName string
	err := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='users';").Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			username TEXT UNIQUE NOT NULL,
			photo_url TEXT DEFAULT ''
		);
		CREATE TABLE IF NOT EXISTS sessions (
			token TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);
		CREATE TABLE IF NOT EXISTS conversations (
			id TEXT PRIMARY KEY,
			name TEXT DEFAULT '',
			photo_url TEXT DEFAULT '',
			is_group BOOLEAN DEFAULT FALSE
		);
		CREATE TABLE IF NOT EXISTS conversation_members (
			conversation_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			PRIMARY KEY (conversation_id, user_id),
			FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);
		CREATE TABLE IF NOT EXISTS messages (
			id TEXT PRIMARY KEY,
			conversation_id TEXT NOT NULL,
			sender_username TEXT NOT NULL,
			content_type TEXT NOT NULL,
			content_value TEXT NOT NULL,
			reply_to_message_id TEXT,
			data_sent DATETIME DEFAULT CURRENT_TIMESTAMP,
			status TEXT DEFAULT 'sent',
			FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE
		);
		CREATE TABLE IF NOT EXISTS reactions (
			id TEXT PRIMARY KEY,
			message_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			reaction_type TEXT NOT NULL,
			FOREIGN KEY (message_id) REFERENCES messages(id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);
		`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
