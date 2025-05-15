package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header missing")
	}
	const prefix = "ApiKey "
	if !strings.HasPrefix(authHeader, prefix) {
		return "", errors.New("authorization header must start with 'ApiKey '")
	}
	return strings.TrimSpace(strings.TrimPrefix(authHeader, prefix)), nil
}
