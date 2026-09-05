package api

import (
	"encoding/json"
	"net/http"
	"regexp"
)

var usernameRegex = regexp.MustCompile(`^.{3,16}$`)

type errorResponse struct {
	Message string `json:"message"`
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, errorResponse{Message: msg})
}

func (rt *_router) auth(w http.ResponseWriter, r *http.Request) (userID string, username string, ok bool) {
	userID, username, err := rt.authenticate(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "Non autorizzato")
		return "", "", false
	}
	return userID, username, true
}
