package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) createConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userID, _, ok := rt.auth(w, r)
	if !ok {
		return
	}

	var payload struct {
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, "JSON non valido")
		return
	}
	if !usernameRegex.MatchString(payload.Username) {
		writeError(w, http.StatusBadRequest, "Username non valido")
		return
	}

	summary, err := rt.db.CreateConversation(userID, payload.Username)
	if err != nil {
		if errors.Is(err, database.ErrUserNotFound) {
			writeError(w, http.StatusNotFound, "Utente non trovato")
			return
		}
		ctx.Logger.WithError(err).Error("createConversation")
		writeError(w, http.StatusInternalServerError, "Errore durante la creazione")
		return
	}

	writeJSON(w, http.StatusCreated, ConversationDetails{
		ID:       summary.ID,
		Username: summary.Username,
		Photo:    summary.Photo,
		IsGroup:  false,
		Messages: []Message{},
	})
}

func (rt *_router) getMyConversations(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userID, _, ok := rt.auth(w, r)
	if !ok {
		return
	}

	convs, err := rt.db.GetUserConversations(userID)
	if err != nil {
		ctx.Logger.WithError(err).Error("getMyConversations")
		writeError(w, http.StatusInternalServerError, "Errore recupero conversazioni")
		return
	}

	out := make([]ConversationDetailsSummary, 0, len(convs))
	for _, c := range convs {
		summary := ConversationDetailsSummary{
			ID:       c.ID,
			Username: c.Username,
			Photo:    c.Photo,
			IsGroup:  c.IsGroup,
		}
		if c.LastMessage != nil {
			summary.LastMessage = &LastMessagePreview{
				Text:           c.LastMessage.Text,
				IsPhoto:        c.LastMessage.IsPhoto,
				DataSent:       c.LastMessage.DataSent.UTC().Format("2006-01-02T15:04:05Z"),
				SenderUsername: c.LastMessage.SenderUsername,
			}
		}
		out = append(out, summary)
	}

	writeJSON(w, http.StatusOK, struct {
		Conversations []ConversationDetailsSummary `json:"conversations"`
	}{Conversations: out})
}

func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userID, _, ok := rt.auth(w, r)
	if !ok {
		return
	}
	convID := ps.ByName("conversationId")

	details, err := rt.db.GetConversationDetails(convID, userID)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			writeError(w, http.StatusNotFound, "Conversazione non trovata")
			return
		}
		ctx.Logger.WithError(err).Error("getConversation")
		writeError(w, http.StatusInternalServerError, "Errore recupero conversazione")
		return
	}

	messages := make([]Message, 0, len(details.Messages))
	for _, m := range details.Messages {
		messages = append(messages, toAPIMessage(m))
	}

	writeJSON(w, http.StatusOK, ConversationDetails{
		ID:       details.ID,
		Username: details.Username,
		Photo:    details.Photo,
		IsGroup:  details.IsGroup,
		Messages: messages,
	})
}
