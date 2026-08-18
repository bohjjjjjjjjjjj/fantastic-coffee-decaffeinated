package api

import (
	"errors"
	"net/http"
	"strings"
)

// authenticate estrae il token dall'header Authorization e recupera userID e username dal DB
func (rt *_router) authenticate(r *http.Request) (string, string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", "", errors.New("missing authorization header")
	}

	// Supporta sia "Bearer <token>" sia solo "<token>"
	token := strings.TrimPrefix(authHeader, "Bearer ")
	token = strings.TrimSpace(token)

	if token == "" {
		return "", "", errors.New("invalid token")
	}

	userID, username, err := rt.db.GetUserByToken(token)
	if err != nil {
		return "", "", errors.New("unauthorized")
	}

	return userID, username, nil
}
