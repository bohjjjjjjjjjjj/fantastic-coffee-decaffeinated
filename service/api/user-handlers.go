package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

// GET /user -> getMyUserInfo
func (rt *_router) getMyUserInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userID, _, ok := rt.auth(w, r)
	if !ok {
		return
	}

	user, err := rt.db.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, database.ErrUserNotFound) {
			writeError(w, http.StatusNotFound, "Utente non trovato")
			return
		}
		ctx.Logger.WithError(err).Error("getMyUserInfo")
		writeError(w, http.StatusInternalServerError, "Errore interno")
		return
	}

	writeJSON(w, http.StatusOK, user)
}

// PUT /user/username -> setMyUserName
func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userID, _, ok := rt.auth(w, r)
	if !ok {
		return
	}

	var payload UserName
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, "JSON non valido")
		return
	}
	if !usernameRegex.MatchString(payload.Username) {
		writeError(w, http.StatusBadRequest, "Username non valido")
		return
	}

	err := rt.db.UpdateUsername(userID, payload.Username)
	if errors.Is(err, database.ErrUsernameTaken) {
		writeError(w, http.StatusConflict, "Username già in uso")
		return
	}
	if err != nil {
		ctx.Logger.WithError(err).Error("setMyUserName")
		writeError(w, http.StatusInternalServerError, "Errore durante l'aggiornamento")
		return
	}

	writeJSON(w, http.StatusOK, payload)
}

// PUT /user/photo -> setMyPhoto
func (rt *_router) setMyPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userID, _, ok := rt.auth(w, r)
	if !ok {
		return
	}

	var payload UserPhoto
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, "JSON non valido")
		return
	}

	if err := rt.db.SetUserPhoto(userID, payload.PhotoURL); err != nil {
		ctx.Logger.WithError(err).Error("setMyPhoto")
		writeError(w, http.StatusInternalServerError, "Errore durante l'aggiornamento")
		return
	}

	writeJSON(w, http.StatusOK, payload)
}

// GET /user/photo -> getMyPhoto
func (rt *_router) getMyPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userID, _, ok := rt.auth(w, r)
	if !ok {
		return
	}

	user, err := rt.db.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, database.ErrUserNotFound) {
			writeError(w, http.StatusNotFound, "Utente non trovato")
			return
		}
		ctx.Logger.WithError(err).Error("getMyPhoto")
		writeError(w, http.StatusInternalServerError, "Errore interno")
		return
	}

	writeJSON(w, http.StatusOK, UserPhoto{PhotoURL: user.PhotoURL})
}

// GET /users -> searchUsers
func (rt *_router) searchUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userID, _, ok := rt.auth(w, r)
	if !ok {
		return
	}

	users, err := rt.db.SearchUsers(r.URL.Query().Get("username"), userID)
	if err != nil {
		ctx.Logger.WithError(err).Error("searchUsers")
		writeError(w, http.StatusInternalServerError, "Errore nella ricerca")
		return
	}

	writeJSON(w, http.StatusOK, struct {
		Users []database.User `json:"users"`
	}{Users: users})
}
