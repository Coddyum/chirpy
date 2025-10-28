package auth

import (
	"fmt"
	"net/http"
	"strings"
)

// Get The Bearer Token inside de Authorization Header
func GetBearerToken(headers http.Header) (string, error) {
	bearerToken := headers.Get("Authorization")
	if bearerToken == "" {
		return "", fmt.Errorf("authorization header missing")
	}

	if len(bearerToken) >= 7 && strings.ToLower(bearerToken[:7]) == "bearer " {
		token := strings.TrimSpace(bearerToken[7:])
		if token == "" {
			return "", fmt.Errorf("empty bearer token")
		}

		return token, nil
	}

	return "", fmt.Errorf("authorization header missing Bearer prefix")
}

// Get The Api Key Token inside de Authorization Header
func GetApiKey(headers http.Header) (string, error) {
	apiKeyToken := headers.Get("Authorization")
	if apiKeyToken == "" {
		return "", fmt.Errorf("authorization header missing")
	}

	if len(apiKeyToken) >= 7 && strings.ToLower(apiKeyToken[:7]) == "apikey " {
		token := strings.TrimSpace(apiKeyToken[7:])
		if token == "" {
			return "", fmt.Errorf("empty apikey token")
		}

		return token, nil
	}

	return "", fmt.Errorf("authorization header missing ApiKey prefix")
}
