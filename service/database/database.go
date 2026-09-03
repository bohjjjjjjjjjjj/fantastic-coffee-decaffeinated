package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// ErrUserNotFound viene restituito quando un utente non esiste.
var ErrUserNotFound = errors.New("user not found")

// ErrNotFound viene restituito quando la risorsa non esiste oppure l'utente non
// ne fa parte. Si usa lo stesso errore nei due casi per rispondere sempre 404 e
// non rivelare l'esistenza di conversazioni/gruppi altrui.
var ErrNotFound = errors.New("resource not found")

// ErrUsernameTaken viene restituito quando lo username richiesto è già in uso.
var ErrUsernameTaken = errors.New("username already taken")

// User rappresenta un utente nel database.
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	PhotoURL string `json:"photoUrl"`
}

// LastMessage è l'anteprima dell'ultimo messaggio di una conversazione.
type LastMessage struct {
	Text           string    `json:"text"`
	DataSent       time.Time `json:"dataSent"`
	SenderUsername string    `json:"senderUsername"`
}

// ConversationSummary rappresenta la sintesi della chat per la lista.
type ConversationSummary struct {
	ID          string
	Username    string
	Photo       string
	IsGroup     bool
	LastMessage *LastMessage
}

// Reaction rappresenta una reazione a un messaggio così come restituita al client.
type Reaction struct {
	ID           string
	ReactionType string
	UserSenderID string
	Username     string
}

// Message rappresenta un messaggio completo di una conversazione.
type Message struct {
	ID               string
	ConversationID   string
	SenderID         string
	SenderUsername   string
	ContentType      string // "text" oppure "photo"
	ContentValue     string
	ReplyToMessageID string
	DataSent         time.Time
	Status           string // calcolato: "sent", "delivered", "read"
	Reactions        []Reaction
}

// ConversationDetails contiene la conversazione con tutti i suoi messaggi.
type ConversationDetails struct {
	ID       string
	Username string
	Photo    string
	IsGroup  bool
	Messages []Message
}

// GroupResponse rappresenta i dettagli di un gruppo.
type GroupResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	PhotoURL string `json:"photoUrl"`
}

// AppDatabase è l'interfaccia che definisce le operazioni sul database.
type AppDatabase interface {
	// Autenticazione e sessioni
	GetOrCreateUser(username string) (userID string, token string, created bool, err error)
	GetUserByToken(token string) (userID string, username string, err error)

	// Profilo utente e ricerca
	GetUserByID(userID string) (User, error)
	UpdateUsername(userID string, newUsername string) error
	SetUserPhoto(userID string, photoURL string) error
	SearchUsers(searchQuery string, excludeUserID string) ([]User, error)

	// Conversazioni
	CreateConversation(currentUserID, targetUsername string) (ConversationSummary, error)
	GetUserConversations(userID string) ([]ConversationSummary, error)
	GetConversationDetails(convID, userID string) (ConversationDetails, error)

	// Messaggi
	SendMessage(convID, senderID, contentType, contentValue, replyToID string) (Message, error)
	ForwardMessage(originalMsgID, targetConvID, senderID string) (Message, error)
	DeleteMessage(msgID, convID, senderID string) error

	// Reazioni
	AddReaction(convID, messageID, userID, reactionType string) (Reaction, error)
	RemoveReaction(convID, messageID, reactionID, userID string) error

	// Gruppi
	CreateGroup(groupName, ownerID string, memberIDs []string) (GroupResponse, error)
	GetGroupDetails(groupID, userID string) (GroupResponse, error)
	GetGroupMembers(groupID, userID string) ([]User, error)
	AddToGroup(groupID, userID string, memberIDs []string) ([]User, error)
	SetGroupName(groupID, newName, userID string) error
	SetGroupPhoto(groupID, photoURL, userID string) error
	LeaveGroup(groupID, userID string) error

	// Controllo stato
	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// schema contiene la struttura completa del database. Tutte le CREATE sono
// idempotenti (IF NOT EXISTS): lo schema viene applicato a ogni avvio.
const schema = `
CREATE TABLE IF NOT EXISTS users (
	id TEXT PRIMARY KEY,
	username TEXT UNIQUE NOT NULL,
	photo_url TEXT NOT NULL DEFAULT ''
);
CREATE TABLE IF NOT EXISTS sessions (
	token TEXT PRIMARY KEY,
	user_id TEXT NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS conversations (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL DEFAULT '',
	photo_url TEXT NOT NULL DEFAULT '',
	is_group BOOLEAN NOT NULL DEFAULT 0
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
	sender_id TEXT NOT NULL,
	content_type TEXT NOT NULL,
	content_value TEXT NOT NULL,
	reply_to_message_id TEXT,
	data_sent DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE,
	FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE,
	FOREIGN KEY (reply_to_message_id) REFERENCES messages(id) ON DELETE SET NULL
);
CREATE TABLE IF NOT EXISTS message_status (
	message_id TEXT NOT NULL,
	user_id TEXT NOT NULL,
	status TEXT NOT NULL,
	PRIMARY KEY (message_id, user_id),
	FOREIGN KEY (message_id) REFERENCES messages(id) ON DELETE CASCADE,
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS reactions (
	id TEXT PRIMARY KEY,
	message_id TEXT NOT NULL,
	user_id TEXT NOT NULL,
	reaction_type TEXT NOT NULL,
	UNIQUE (message_id, user_id),
	FOREIGN KEY (message_id) REFERENCES messages(id) ON DELETE CASCADE,
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
`

// New crea e inizializza una nuova istanza di AppDatabase su una connessione SQLite.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("a database is required")
	}

	if _, err := db.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		return nil, fmt.Errorf("error enabling foreign keys: %w", err)
	}

	if _, err := db.Exec(schema); err != nil {
		return nil, fmt.Errorf("error creating database structure: %w", err)
	}

	return &appdbimpl{c: db}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}

// isMember indica se l'utente fa parte della conversazione indicata.
func (db *appdbimpl) isMember(convID, userID string) (bool, error) {
	var ok bool
	err := db.c.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM conversation_members WHERE conversation_id = ? AND user_id = ?)`,
		convID, userID,
	).Scan(&ok)
	return ok, err
}

// conversationView restituisce nome e foto della conversazione dal punto di vista
// dell'utente: per i gruppi usa nome/foto del gruppo, per le chat dirette quelli
// dell'altro partecipante.
func (db *appdbimpl) conversationView(convID, userID string) (title string, photo string, isGroup bool, err error) {
	err = db.c.QueryRow(`
		SELECT
			c.is_group,
			CASE WHEN c.is_group THEN c.name ELSE COALESCE((
				SELECT u.username FROM conversation_members cm
				JOIN users u ON u.id = cm.user_id
				WHERE cm.conversation_id = c.id AND cm.user_id != ?
				LIMIT 1
			), '') END,
			CASE WHEN c.is_group THEN c.photo_url ELSE COALESCE((
				SELECT u.photo_url FROM conversation_members cm
				JOIN users u ON u.id = cm.user_id
				WHERE cm.conversation_id = c.id AND cm.user_id != ?
				LIMIT 1
			), '') END
		FROM conversations c
		WHERE c.id = ?`,
		userID, userID, convID,
	).Scan(&isGroup, &title, &photo)
	return title, photo, isGroup, err
}
