package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var ErrUserNotFound = errors.New("user not found")

var ErrNotFound = errors.New("resource not found")

var ErrUsernameTaken = errors.New("username already taken")

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	PhotoURL string `json:"photoUrl"`
}

type LastMessage struct {
	Text           string    `json:"text"`
	IsPhoto        bool      `json:"isPhoto"`
	DataSent       time.Time `json:"dataSent"`
	SenderUsername string    `json:"senderUsername"`
}

type ConversationSummary struct {
	ID          string
	Username    string
	Photo       string
	IsGroup     bool
	LastMessage *LastMessage
}

type Reaction struct {
	ID           string
	ReactionType string
	UserSenderID string
	Username     string
}

type Message struct {
	ID               string
	ConversationID   string
	SenderID         string
	SenderUsername   string
	SenderPhotoURL   string
	ContentText      string
	ContentPhoto     string
	ReplyToMessageID string
	IsForwarded      bool
	DataSent         time.Time
	Status           string
	Reactions        []Reaction
}

type ConversationDetails struct {
	ID       string
	Username string
	Photo    string
	IsGroup  bool
	Messages []Message
}

type GroupResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	PhotoURL string `json:"photoUrl"`
}

type AppDatabase interface {
	GetOrCreateUser(username string) (userID string, token string, created bool, err error)
	GetUserByToken(token string) (userID string, username string, err error)

	GetUserByID(userID string) (User, error)
	UpdateUsername(userID string, newUsername string) error
	SetUserPhoto(userID string, photoURL string) error
	SearchUsers(searchQuery string, excludeUserID string) ([]User, error)

	CreateConversation(currentUserID, targetUsername string) (ConversationSummary, error)
	GetUserConversations(userID string) ([]ConversationSummary, error)
	GetConversationDetails(convID, userID string) (ConversationDetails, error)

	SendMessage(convID, senderID, contentText, contentPhoto, replyToID string) (Message, error)
	ForwardMessage(originalMsgID, targetConvID, senderID string) (Message, error)
	DeleteMessage(msgID, convID, senderID string) error

	AddReaction(convID, messageID, userID, reactionType string) (Reaction, error)
	RemoveReaction(convID, messageID, reactionID, userID string) error

	CreateGroup(groupName, ownerID string, memberIDs []string) (GroupResponse, error)
	GetGroupDetails(groupID, userID string) (GroupResponse, error)
	GetGroupMembers(groupID, userID string) ([]User, error)
	AddToGroup(groupID, userID string, memberIDs []string) ([]User, error)
	SetGroupName(groupID, newName, userID string) error
	SetGroupPhoto(groupID, photoURL, userID string) error
	LeaveGroup(groupID, userID string) error

	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

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
	content_text TEXT NOT NULL DEFAULT '',
	content_photo TEXT NOT NULL DEFAULT '',
	reply_to_message_id TEXT,
	is_forwarded BOOLEAN NOT NULL DEFAULT 0,
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

	// === SEED UTENTI PREDEFINITI ===
	seedQuery := `
	INSERT OR IGNORE INTO users (id, username, photo_url) VALUES 
		('usr_hiroto', 'hiroto', ''),
		('usr_sakuragi', 'sakuragi', ''),
		('usr_naruto', 'naruto', ''),
		('usr_luffy', 'luffy', ''),
		('usr_tohru', 'tohru', ''),
		('usr_meow', 'meow', '');

	INSERT OR IGNORE INTO sessions (token, user_id) VALUES 
		('sess_hiroto', 'usr_hiroto'),
		('sess_sakuragi', 'usr_sakuragi'),
		('sess_naruto', 'usr_naruto'),
		('sess_luffy', 'usr_luffy'),
		('sess_tohru', 'usr_tohru'),
		('sess_meow', 'usr_meow');
	`
	if _, err := db.Exec(seedQuery); err != nil {
		return nil, fmt.Errorf("error seeding database: %w", err)
	}

	return &appdbimpl{c: db}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}

func (db *appdbimpl) isMember(convID, userID string) (bool, error) {
	var ok bool
	err := db.c.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM conversation_members WHERE conversation_id = ? AND user_id = ?)`,
		convID, userID,
	).Scan(&ok)
	return ok, err
}

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
