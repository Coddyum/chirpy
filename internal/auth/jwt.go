package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type CustomClaims struct {
	jwt.RegisteredClaims
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expriresIn time.Duration) (string, error) {

	claims := CustomClaims{
		jwt.RegisteredClaims{
			Issuer:    "chirpy",
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expriresIn)),
			Subject:   userID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		fmt.Println("Token can't be signed")
		return "", err
	}

	return signedToken, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	}, jwt.WithLeeway(5*time.Second))
	if err != nil {
		return uuid.Nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return uuid.Nil, fmt.Errorf("invalid token claims")
	}

	if claims.ExpiresAt != nil && time.Now().After(claims.ExpiresAt.Time) {
		return uuid.Nil, fmt.Errorf("token expired")
	}

	sub, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, fmt.Errorf("cannot parse subject claim: %w", err)
	}

	return sub, nil
}
