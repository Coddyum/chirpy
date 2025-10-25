package auth

import (
	"fmt"
	"net/http"
	"strings"
)

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
