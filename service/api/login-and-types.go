package api

import (
	"encoding/json"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

type loginRequest struct {
	Name string `json:"name"`
}

type loginResponse struct {
	Identifier string `json:"identifier"`
}

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var payload loginRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, "JSON non valido")
		return
	}

	if !usernameRegex.MatchString(payload.Name) {
		writeError(w, http.StatusBadRequest, "Username non valido (da 3 a 16 caratteri)")
		return
	}

	_, token, _, err := rt.db.GetOrCreateUser(payload.Name)
	if err != nil {
		ctx.Logger.WithError(err).Error("doLogin: error creating user")
		writeError(w, http.StatusInternalServerError, "Errore interno durante il login")
		return
	}

	writeJSON(w, http.StatusCreated, loginResponse{Identifier: token})
}
