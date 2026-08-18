package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

// POST /conversations -> createConversation
func (rt *_router) createConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	userID, _, err := rt.authenticate(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(errorResponse{Message: "Non autorizzato"})
		return
	}

	var payload struct {
		Username string `json:"username"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse{Message: "JSON non valido"})
		return
	}

	summary, err := rt.db.CreateConversation(userID, payload.Username)
	if err != nil {
		if errors.Is(err, database.ErrUserNotFound) {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(errorResponse{Message: "Utente non trovato"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errorResponse{Message: "Errore durante la creazione"})
		return
	}

	response := struct {
		ID       string        `json:"id"`
		Username string        `json:"username"`
		Photo    string        `json:"photo"`
		Messages []interface{} `json:"messages"`
	}{
		ID:       summary.ID,
		Username: summary.Username,
		Photo:    summary.Photo,
		Messages: []interface{}{},
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(response)
}

// GET /conversations -> getMyConversations
func (rt *_router) getMyConversations(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	userID, _, err := rt.authenticate(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(errorResponse{Message: "Non autorizzato"})
		return
	}

	convs, err := rt.db.GetUserConversations(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errorResponse{Message: "Errore recupero conversazioni"})
		return
	}

	response := struct {
		Conversations []database.ConversationSummary `json:"conversations"`
	}{
		Conversations: convs,
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}
