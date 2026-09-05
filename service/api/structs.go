package api

import "git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"

type UserName struct {
	Username string `json:"username"`
}

type UserPhoto struct {
	PhotoURL string `json:"photoUrl"`
}

type Reaction struct {
	ID           string `json:"id"`
	ReactionType string `json:"reactionType"`
	UserSenderID string `json:"userSenderId"`
	Username     string `json:"username,omitempty"`
}

type StatusMessage struct {
	Value string `json:"value"`
}

type MessageContent struct {
	Text     string `json:"text,omitempty"`
	PhotoURL string `json:"photoUrl,omitempty"`
}

type Message struct {
	ID               string         `json:"id"`
	ConversationID   string         `json:"conversationId"`
	SenderID         string         `json:"senderId"`
	SenderUsername   string         `json:"senderUsername"`
	SenderPhotoURL   string         `json:"senderPhotoUrl"`
	DataSent         string         `json:"dataSent"`
	Status           StatusMessage  `json:"status"`
	Content          MessageContent `json:"content"`
	ReplyToMessageID string         `json:"replyToMessageId,omitempty"`
	IsForwarded      bool           `json:"isForwarded"`
	Reactions        []Reaction     `json:"reaction"`
}

type NewMessageRequest struct {
	Content          MessageContent `json:"content"`
	ReplyToMessageID string         `json:"replyToMessageId,omitempty"`
}

type MessageForward struct {
	TargetConversationID string `json:"targetConversationId"`
}

type LastMessagePreview struct {
	Text           string `json:"text"`
	IsPhoto        bool   `json:"isPhoto"`
	DataSent       string `json:"dataSent"`
	SenderUsername string `json:"senderUsername"`
}

type ConversationDetailsSummary struct {
	ID          string              `json:"id"`
	Username    string              `json:"username"`
	Photo       string              `json:"photo"`
	IsGroup     bool                `json:"isGroup"`
	LastMessage *LastMessagePreview `json:"lastMessage,omitempty"`
}

type ConversationDetails struct {
	ID       string    `json:"id"`
	Username string    `json:"username"`
	Photo    string    `json:"photo"`
	IsGroup  bool      `json:"isGroup"`
	Messages []Message `json:"messages"`
}

type Group struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	PhotoURL string `json:"photoUrl"`
}

type GroupName struct {
	Name string `json:"name"`
}

type GroupPhoto struct {
	PhotoURL string `json:"photoUrl"`
}

type MemberIDList struct {
	MemberIDs []string `json:"memberIds"`
}

type GroupMemberListResponse struct {
	Members []database.User `json:"members"`
}

func toAPIMessage(m database.Message) Message {
	content := MessageContent{Text: m.ContentText, PhotoURL: m.ContentPhoto}

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
		SenderPhotoURL:   m.SenderPhotoURL,
		DataSent:         m.DataSent.UTC().Format("2006-01-02T15:04:05Z"),
		Status:           StatusMessage{Value: m.Status},
		Content:          content,
		ReplyToMessageID: m.ReplyToMessageID,
		IsForwarded:      m.IsForwarded,
		Reactions:        reactions,
	}
}
