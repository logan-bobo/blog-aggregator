package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(header http.Header) (string, error) {
	authHeader := header.Get("Authorization")

	if authHeader == "" {
		return "", errors.New("No authorization header included")
	}

	splitAuth := strings.Split(authHeader, " ")

	if len(splitAuth) < 2 || splitAuth[0] != "ApiKey" {
		return "", errors.New("Malformed authorization header")
	}

	return splitAuth[1], nil
}
