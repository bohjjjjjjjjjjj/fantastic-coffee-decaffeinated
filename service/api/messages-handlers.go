package api

import (
	"encoding/json"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// POST /conversations/:conversationId/messages -> sendMessage
func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	convID := ps.ByName("conversationId")

	_, username, err := rt.authenticate(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(errorResponse{Message: "Non autorizzato"})
		return
	}

	var payload struct {
		Content struct {
			Text     string `json:"text"`
			PhotoURL string `json:"photoUrl"`
		} `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse{Message: "JSON non valido"})
		return
	}

	cType := "text"
	cVal := payload.Content.Text
	if payload.Content.PhotoURL != "" {
		cType = "photo"
		cVal = payload.Content.PhotoURL
	}

	if cVal == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse{Message: "Contenuto vuoto"})
		return
	}

	msg, err := rt.db.SendMessage(convID, username, cType, cVal)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errorResponse{Message: "Errore invio messaggio"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(msg)
}

// GET /conversations/:conversationId -> getConversation
func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	convID := ps.ByName("conversationId")

	userID, _, err := rt.authenticate(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(errorResponse{Message: "Non autorizzato"})
		return
	}

	details, err := rt.db.GetConversationDetails(convID, userID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(errorResponse{Message: "Conversazione non trovata"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(details)
}

// DELETE /conversations/:conversationId/messages/:messageId -> deleteMessage
func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	convID := ps.ByName("conversationId")
	msgID := ps.ByName("messageId")

	_, username, err := rt.authenticate(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(errorResponse{Message: "Non autorizzato"})
		return
	}

	err = rt.db.DeleteMessage(msgID, convID, username)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(errorResponse{Message: "Impossibile eliminare il messaggio"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
