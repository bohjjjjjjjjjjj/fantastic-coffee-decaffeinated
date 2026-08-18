package api

// User rappresenta le informazioni di base di un utente (schema User)
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	PhotoURL string `json:"photoUrl,omitempty"`
}

// UserName viene usato per il payload di aggiornamento del nome utente (schema UserName)
type UserName struct {
	Username string `json:"username"`
}

// UserPhoto viene usato per il payload di aggiornamento della foto profilo (schema UserPhoto)
type UserPhoto struct {
	PhotoURL string `json:"photoUrl"`
}

// Reaction rappresenta una reazione ad un messaggio (schema Reaction)
type Reaction struct {
	ReactionType string `json:"reactionType"`
	UserSenderID string `json:"userSenderId"`
}

// NewReaction è il payload inviato dal client per aggiungere una reazione (schema NewReaction)
type NewReaction struct {
	ReactionType string `json:"reactionType"`
}

// StatusMessage descrive lo stato del messaggio (schema StatusMessage)
type StatusMessage struct {
	Value string `json:"value"` // "sent", "delivered", "read", "failed"
}

// TextMessageContent rappresenta un messaggio di solo testo (schema TextMessageContent)
type TextMessageContent struct {
	Text string `json:"text"`
}

// PhotoMessageContent rappresenta un messaggio contenente un'immagine (schema PhotoMessageContent)
type PhotoMessageContent struct {
	PhotoURL string `json:"photoUrl"`
}

// MessageContent wrapper per gestire sia messaggi di testo che di foto
type MessageContent struct {
	Text     string `json:"text,omitempty"`
	PhotoURL string `json:"photoUrl,omitempty"`
}

// Message rappresenta un messaggio completo (schema Message)
type Message struct {
	ID             string         `json:"id"`
	SenderUsername string         `json:"senderUsername"`
	DataSent       string         `json:"dataSent"`
	Status         StatusMessage  `json:"status"`
	Content        MessageContent `json:"content"`
	Reactions      []Reaction     `json:"reaction,omitempty"`
}

// NewMessageRequest è il payload per l'invio di un messaggio (schema NewMessageRequest)
type NewMessageRequest struct {
	Content          MessageContent `json:"content"`
	ReplyToMessageID string         `json:"replyToMessageId,omitempty"`
}

// MessageForward è il payload per l'inoltro di un messaggio (schema MessageForward)
type MessageForward struct {
	TargetConversationID string `json:"targetConversationId"`
}

// ConversationDetailsSummary è la sintesi per l'elenco chat (schema ConversationDetailsSummary)
type ConversationDetailsSummary struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Photo    string `json:"photo"`
}

// ConversationDetails contiene la chat completa con i messaggi (schema ConversationDetails)
type ConversationDetails struct {
	ID       string    `json:"id"`
	Username string    `json:"username"`
	Photo    string    `json:"photo"`
	Messages []Message `json:"messages"`
}

// Group rappresenta i dettagli di un gruppo (schema Group)
type Group struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	PhotoURL string `json:"photoUrl"`
}

// GroupName payload per aggiornare il nome del gruppo (schema GroupName)
type GroupName struct {
	Name string `json:"name"`
}

// GroupPhoto payload per aggiornare la foto del gruppo (schema GroupPhoto)
type GroupPhoto struct {
	PhotoURL string `json:"photoUrl"`
}

// MemberIDList payload per aggiungere membri al gruppo (schema MemberIdList)
type MemberIDList struct {
	MemberIDs []string `json:"memberIds"`
}

// GroupMemberListResponse contiene l'elenco aggiornato dei membri di un gruppo
type GroupMemberListResponse struct {
	Members []User `json:"members"`
}

// LoginRequest payload per la richiesta di login
type LoginRequest struct {
	Name string `json:"name"`
}

// LoginResponse risposta restitutita dopo il login
type LoginResponse struct {
	Identifier string `json:"identifier"`
}

// ErrorResponse formato di errore standard (schema ErrorResponse)
type ErrorResponse struct {
	Message string `json:"message"`
}
