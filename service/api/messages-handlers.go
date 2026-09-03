package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

// POST /conversations/:conversationId/messages -> sendMessage
func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userID, _, ok := rt.auth(w, r)
	if !ok {
		return
	}
	convID := ps.ByName("conversationId")

	var payload NewMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, "JSON non valido")
		return
	}

	cType := "text"
	cVal := payload.Content.Text
	if payload.Content.PhotoURL != "" {
		cType = "photo"
		cVal = payload.Content.PhotoURL
	}
	if cVal == "" {
		writeError(w, http.StatusBadRequest, "Contenuto vuoto")
		return
	}

	msg, err := rt.db.SendMessage(convID, userID, cType, cVal, payload.ReplyToMessageID)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			writeError(w, http.StatusNotFound, "Conversazione non trovata")
			return
		}
		ctx.Logger.WithError(err).Error("sendMessage")
		writeError(w, http.StatusInternalServerError, "Errore invio messaggio")
		return
	}

	writeJSON(w, http.StatusCreated, toAPIMessage(msg))
}

// POST /conversations/:conversationId/messages/:messageId/forward -> forwardMessage
func (rt *_router) forwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userID, _, ok := rt.auth(w, r)
	if !ok {
		return
	}
	msgID := ps.ByName("messageId")

	var payload MessageForward
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.TargetConversationID == "" {
		writeError(w, http.StatusBadRequest, "Destinazione non valida")
		return
	}

	newMsg, err := rt.db.ForwardMessage(msgID, payload.TargetConversationID, userID)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			writeError(w, http.StatusNotFound, "Messaggio o conversazione non trovati")
			return
		}
		ctx.Logger.WithError(err).Error("forwardMessage")
		writeError(w, http.StatusInternalServerError, "Errore inoltro messaggio")
		return
	}

	writeJSON(w, http.StatusCreated, toAPIMessage(newMsg))
}

// DELETE /conversations/:conversationId/messages/:messageId -> deleteMessage
func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userID, _, ok := rt.auth(w, r)
	if !ok {
		return
	}
	convID := ps.ByName("conversationId")
	msgID := ps.ByName("messageId")

	err := rt.db.DeleteMessage(msgID, convID, userID)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			writeError(w, http.StatusNotFound, "Messaggio non trovato")
			return
		}
		ctx.Logger.WithError(err).Error("deleteMessage")
		writeError(w, http.StatusInternalServerError, "Errore eliminazione messaggio")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
