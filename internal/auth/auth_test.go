package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeAndValidateJWT(t *testing.T) {
	secret := "supersecret"
	userID := uuid.New()

	token, err := MakeJWT(userID, secret, time.Minute)
	if err != nil {
		t.Fatalf("failed to make jwt: %v", err)
	}

	parsedID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("failed to validate jwt: %v", err)
	}

	if parsedID != userID {
		t.Errorf("expected userID %v, got %v", userID, parsedID)
	}
}

func TestExpiredJWT(t *testing.T) {
	secret := "supersecret"
	userID := uuid.New()

	token, err := MakeJWT(userID, secret, -time.Second)
	if err != nil {
		t.Fatalf("failed to make jwt: %v", err)
	}

	_, err = ValidateJWT(token, secret)
	if err == nil {
		t.Errorf("expected error for expired token, got nil")
	}
}

func TestInvalidSecretJWT(t *testing.T) {
	userID := uuid.New()
	token, err := MakeJWT(userID, "secret1", time.Minute)
	if err != nil {
		t.Fatalf("failed to make jwt: %v", err)
	}

	_, err = ValidateJWT(token, "wrongsecret")
	if err == nil {
		t.Errorf("expected error for invalid secret, got nil")
	}
}
