package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) createGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userID, _, ok := rt.auth(w, r)
	if !ok {
		return
	}

	var payload struct {
		Name      string   `json:"name"`
		MemberIDs []string `json:"memberIds"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.Name == "" {
		writeError(w, http.StatusBadRequest, "Nome gruppo non valido")
		return
	}

	group, err := rt.db.CreateGroup(payload.Name, userID, payload.MemberIDs)
	if err != nil {
		ctx.Logger.WithError(err).Error("createGroup")
		writeError(w, http.StatusInternalServerError, "Errore creazione gruppo")
		return
	}

	writeJSON(w, http.StatusCreated, Group{ID: group.ID, Name: group.Name, PhotoURL: group.PhotoURL})
}

func (rt *_router) getGroupDetails(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userID, _, ok := rt.auth(w, r)
	if !ok {
		return
	}
	groupID := ps.ByName("groupId")

	group, err := rt.db.GetGroupDetails(groupID, userID)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			writeError(w, http.StatusNotFound, "Gruppo non trovato")
			return
		}
		ctx.Logger.WithError(err).Error("getGroupDetails")
		writeError(w, http.StatusInternalServerError, "Errore recupero gruppo")
		return
	}

	writeJSON(w, http.StatusOK, Group{ID: group.ID, Name: group.Name, PhotoURL: group.PhotoURL})
}

func (rt *_router) getGroupMembers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userID, _, ok := rt.auth(w, r)
	if !ok {
		return
	}
	groupID := ps.ByName("groupId")

	members, err := rt.db.GetGroupMembers(groupID, userID)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			writeError(w, http.StatusNotFound, "Gruppo non trovato")
			return
		}
		ctx.Logger.WithError(err).Error("getGroupMembers")
		writeError(w, http.StatusInternalServerError, "Errore recupero membri")
		return
	}

	writeJSON(w, http.StatusOK, GroupMemberListResponse{Members: members})
}

func (rt *_router) addToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userID, _, ok := rt.auth(w, r)
	if !ok {
		return
	}
	groupID := ps.ByName("groupId")

	var payload MemberIDList
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || len(payload.MemberIDs) == 0 {
		writeError(w, http.StatusBadRequest, "Elenco membri non valido")
		return
	}

	members, err := rt.db.AddToGroup(groupID, userID, payload.MemberIDs)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			writeError(w, http.StatusNotFound, "Gruppo non trovato")
			return
		}
		ctx.Logger.WithError(err).Error("addToGroup")
		writeError(w, http.StatusInternalServerError, "Errore aggiunta membri")
		return
	}

	writeJSON(w, http.StatusOK, GroupMemberListResponse{Members: members})
}

func (rt *_router) setGroupName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userID, _, ok := rt.auth(w, r)
	if !ok {
		return
	}
	groupID := ps.ByName("groupId")

	var payload GroupName
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.Name == "" {
		writeError(w, http.StatusBadRequest, "Nome gruppo non valido")
		return
	}

	err := rt.db.SetGroupName(groupID, payload.Name, userID)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			writeError(w, http.StatusNotFound, "Gruppo non trovato")
			return
		}
		ctx.Logger.WithError(err).Error("setGroupName")
		writeError(w, http.StatusInternalServerError, "Errore aggiornamento nome")
		return
	}

	writeJSON(w, http.StatusOK, payload)
}

func (rt *_router) setGroupPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userID, _, ok := rt.auth(w, r)
	if !ok {
		return
	}
	groupID := ps.ByName("groupId")

	var payload GroupPhoto
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, "JSON non valido")
		return
	}

	err := rt.db.SetGroupPhoto(groupID, payload.PhotoURL, userID)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			writeError(w, http.StatusNotFound, "Gruppo non trovato")
			return
		}
		ctx.Logger.WithError(err).Error("setGroupPhoto")
		writeError(w, http.StatusInternalServerError, "Errore aggiornamento foto")
		return
	}

	writeJSON(w, http.StatusOK, payload)
}

func (rt *_router) leaveGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userID, _, ok := rt.auth(w, r)
	if !ok {
		return
	}
	groupID := ps.ByName("groupId")

	err := rt.db.LeaveGroup(groupID, userID)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			writeError(w, http.StatusNotFound, "Gruppo non trovato")
			return
		}
		ctx.Logger.WithError(err).Error("leaveGroup")
		writeError(w, http.StatusInternalServerError, "Errore uscita dal gruppo")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
