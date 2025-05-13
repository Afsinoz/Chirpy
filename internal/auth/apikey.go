package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	authValues := headers["Authorization"]

	prefix := "ApiKey"

	if len(authValues) == 0 {
		return "", errors.New("authorization header not found")
	}

	authHeader := authValues[0]

	if !strings.HasPrefix(authHeader, prefix) {
		return "", errors.New("authorization header is not in expected format")
	}

	apiKey := authHeader[len(prefix)+1:]

	if apiKey == "" {
		return "", errors.New("apikey is not found")
	}

	return apiKey, nil
}
