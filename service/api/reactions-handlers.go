package api

import (
	"encoding/json"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// POST /conversations/:conversationId/messages/:messageId/reactions -> commentMessage
func (rt *_router) commentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	msgID := ps.ByName("messageId")

	userID, _, err := rt.authenticate(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(errorResponse{Message: "Non autorizzato"})
		return
	}

	var payload struct {
		ReactionType string `json:"reactionType"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.ReactionType == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse{Message: "Reazione non valida"})
		return
	}

	reaction, err := rt.db.AddReaction(msgID, userID, payload.ReactionType)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(errorResponse{Message: "Messaggio non trovato"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(reaction)
}

// DELETE /conversations/:conversationId/messages/:messageId/reactions/:reactionId -> uncommentMessage
func (rt *_router) uncommentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	msgID := ps.ByName("messageId")
	reactionID := ps.ByName("reactionId")

	userID, _, err := rt.authenticate(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(errorResponse{Message: "Non autorizzato"})
		return
	}

	err = rt.db.RemoveReaction(msgID, reactionID, userID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(errorResponse{Message: "Reazione non trovata"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// POST /conversations/:conversationId/messages/:messageId/forward -> forwardMessage
func (rt *_router) forwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	msgID := ps.ByName("messageId")

	_, username, err := rt.authenticate(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(errorResponse{Message: "Non autorizzato"})
		return
	}

	var payload struct {
		TargetConversationID string `json:"targetConversationId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.TargetConversationID == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse{Message: "Destinazione non valida"})
		return
	}

	newMsg, err := rt.db.ForwardMessage(msgID, payload.TargetConversationID, username)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(errorResponse{Message: "Messaggio o conversazione non trovati"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(newMsg)
}
