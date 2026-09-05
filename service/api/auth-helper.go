package api

import (
	"errors"
	"net/http"
	"strings"
)

func (rt *_router) authenticate(r *http.Request) (string, string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", "", errors.New("missing authorization header")
	}

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
