package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) commentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userID, _, ok := rt.auth(w, r)
	if !ok {
		return
	}
	convID := ps.ByName("conversationId")
	msgID := ps.ByName("messageId")

	var payload struct {
		ReactionType string `json:"reactionType"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.ReactionType == "" {
		writeError(w, http.StatusBadRequest, "Reazione non valida")
		return
	}

	reaction, err := rt.db.AddReaction(convID, msgID, userID, payload.ReactionType)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			writeError(w, http.StatusNotFound, "Messaggio non trovato")
			return
		}
		ctx.Logger.WithError(err).Error("commentMessage")
		writeError(w, http.StatusInternalServerError, "Errore aggiunta reazione")
		return
	}

	writeJSON(w, http.StatusCreated, Reaction{
		ID:           reaction.ID,
		ReactionType: reaction.ReactionType,
		UserSenderID: reaction.UserSenderID,
		Username:     reaction.Username,
	})
}

func (rt *_router) uncommentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userID, _, ok := rt.auth(w, r)
	if !ok {
		return
	}
	convID := ps.ByName("conversationId")
	msgID := ps.ByName("messageId")
	reactionID := ps.ByName("reactionId")

	err := rt.db.RemoveReaction(convID, msgID, reactionID, userID)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			writeError(w, http.StatusNotFound, "Reazione non trovata")
			return
		}
		ctx.Logger.WithError(err).Error("uncommentMessage")
		writeError(w, http.StatusInternalServerError, "Errore rimozione reazione")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
