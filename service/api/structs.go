package api

import "git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"

// UserName è il payload per l'aggiornamento del nome utente (schema UserName).
type UserName struct {
	Username string `json:"username"`
}

// UserPhoto è il payload / risposta per la foto profilo (schema UserPhoto).
type UserPhoto struct {
	PhotoURL string `json:"photoUrl"`
}

// Reaction rappresenta una reazione a un messaggio (schema Reaction).
type Reaction struct {
	ID           string `json:"id"`
	ReactionType string `json:"reactionType"`
	UserSenderID string `json:"userSenderId"`
	Username     string `json:"username,omitempty"`
}

// StatusMessage descrive lo stato del messaggio (schema StatusMessage).
type StatusMessage struct {
	Value string `json:"value"` // "sent", "delivered", "read", "failed"
}

// MessageContent gestisce sia i messaggi di testo sia quelli con foto.
type MessageContent struct {
	Text     string `json:"text,omitempty"`
	PhotoURL string `json:"photoUrl,omitempty"`
}

// Message rappresenta un messaggio completo (schema Message).
type Message struct {
	ID               string         `json:"id"`
	ConversationID   string         `json:"conversationId"`
	SenderID         string         `json:"senderId"`
	SenderUsername   string         `json:"senderUsername"`
	DataSent         string         `json:"dataSent"`
	Status           StatusMessage  `json:"status"`
	Content          MessageContent `json:"content"`
	ReplyToMessageID string         `json:"replyToMessageId,omitempty"`
	IsForwarded      bool           `json:"isForwarded"`
	Reactions        []Reaction     `json:"reaction"`
}

// NewMessageRequest è il payload per l'invio di un messaggio (schema NewMessageRequest).
type NewMessageRequest struct {
	Content          MessageContent `json:"content"`
	ReplyToMessageID string         `json:"replyToMessageId,omitempty"`
}

// MessageForward è il payload per l'inoltro di un messaggio (schema MessageForward).
type MessageForward struct {
	TargetConversationID string `json:"targetConversationId"`
}

// LastMessagePreview è l'anteprima dell'ultimo messaggio nella lista chat.
type LastMessagePreview struct {
	Text           string `json:"text"`
	DataSent       string `json:"dataSent"`
	SenderUsername string `json:"senderUsername"`
}

// ConversationDetailsSummary è la sintesi per l'elenco chat (schema ConversationDetailsSummary).
type ConversationDetailsSummary struct {
	ID          string              `json:"id"`
	Username    string              `json:"username"`
	Photo       string              `json:"photo"`
	IsGroup     bool                `json:"isGroup"`
	LastMessage *LastMessagePreview `json:"lastMessage,omitempty"`
}

// ConversationDetails contiene la chat completa con i messaggi (schema ConversationDetails).
type ConversationDetails struct {
	ID       string    `json:"id"`
	Username string    `json:"username"`
	Photo    string    `json:"photo"`
	IsGroup  bool      `json:"isGroup"`
	Messages []Message `json:"messages"`
}

// Group rappresenta i dettagli di un gruppo (schema Group).
type Group struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	PhotoURL string `json:"photoUrl"`
}

// GroupName è il payload per aggiornare il nome del gruppo (schema GroupName).
type GroupName struct {
	Name string `json:"name"`
}

// GroupPhoto è il payload per aggiornare la foto del gruppo (schema GroupPhoto).
type GroupPhoto struct {
	PhotoURL string `json:"photoUrl"`
}

// MemberIDList è il payload per aggiungere membri al gruppo (schema MemberIdList).
type MemberIDList struct {
	MemberIDs []string `json:"memberIds"`
}

// GroupMemberListResponse contiene l'elenco aggiornato dei membri di un gruppo.
type GroupMemberListResponse struct {
	Members []database.User `json:"members"`
}

// toAPIMessage converte un messaggio del database nella rappresentazione API.
func toAPIMessage(m database.Message) Message {
	content := MessageContent{}
	if m.ContentType == "photo" {
		content.PhotoURL = m.ContentValue
	} else {
		content.Text = m.ContentValue
	}

	reactions := make([]Reaction, 0, len(m.Reactions))
	for _, r := range m.Reactions {
		reactions = append(reactions, Reaction{
			ID:           r.ID,
			ReactionType: r.ReactionType,
			UserSenderID: r.UserSenderID,
			Username:     r.Username,
		})
	}

	return Message{
		ID:               m.ID,
		ConversationID:   m.ConversationID,
		SenderID:         m.SenderID,
		SenderUsername:   m.SenderUsername,
		DataSent:         m.DataSent.UTC().Format("2006-01-02T15:04:05Z"),
		Status:           StatusMessage{Value: m.Status},
		Content:          content,
		ReplyToMessageID: m.ReplyToMessageID,
		IsForwarded:      m.IsForwarded,
		Reactions:        reactions,
	}
}
