package api

import (
	"encoding/json"
	"net/http"
	"regexp"
)

// usernameRegex valida gli username: 3-30 caratteri alfanumerici o underscore.
var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,30}$`)

// errorResponse è il formato standard per le risposte di errore (schema ErrorResponse).
type errorResponse struct {
	Message string `json:"message"`
}

// writeJSON serializza v come JSON con lo status indicato.
func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// writeError scrive una risposta di errore JSON con lo status indicato.
func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, errorResponse{Message: msg})
}

// auth verifica il token Bearer e restituisce userID e username. In caso di
// fallimento scrive già la risposta 401 e ritorna ok = false.
func (rt *_router) auth(w http.ResponseWriter, r *http.Request) (userID string, username string, ok bool) {
	userID, username, err := rt.authenticate(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "Non autorizzato")
		return "", "", false
	}
	return userID, username, true
}
