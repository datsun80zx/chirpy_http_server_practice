package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("no authorization header found")
	}
	apiKey, found := strings.CutPrefix(authHeader, "ApiKey ")
	if !found {
		return "", fmt.Errorf("no apikey found")
	}
	return apiKey, nil
}
