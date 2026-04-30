package utils

import (
	"regexp"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHashAndCheckPassword(t *testing.T) {
	t.Parallel()

	password := "s3cr3tP@ssw0rd"

	hash1, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}
	if hash1 == password {
		t.Fatalf("expected hash to differ from password")
	}

	hash2, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned error on second call: %v", err)
	}
	if hash2 == password {
		t.Fatalf("expected second hash to differ from password")
	}
	if hash1 == hash2 {
		t.Fatalf("expected two hashes for same password to differ due to salt")
	}

	if !CheckPasswordHash(password, hash1) {
		t.Fatalf("CheckPasswordHash returned false for correct password/hash1")
	}
	if !CheckPasswordHash(password, hash2) {
		t.Fatalf("CheckPasswordHash returned false for correct password/hash2")
	}
	if CheckPasswordHash("wrongpassword", hash1) {
		t.Fatalf("CheckPasswordHash returned true for incorrect password")
	}
}

func TestGenerateSessionID(t *testing.T) {
	t.Parallel()

	id1, err := GenerateSessionID()
	if err != nil {
		t.Fatalf("GenerateSessionID returned error: %v", err)
	}
	if id1 == "" {
		t.Fatalf("GenerateSessionID returned empty id")
	}
	if _, err := uuid.Parse(id1); err != nil {
		t.Fatalf("GenerateSessionID returned invalid UUID: %v", err)
	}

	id2, err := GenerateSessionID()
	if err != nil {
		t.Fatalf("GenerateSessionID returned error on second call: %v", err)
	}
	if id1 == id2 {
		t.Fatalf("expected two generated session IDs to differ")
	}
}

func TestGenerateRandomToken(t *testing.T) {
	t.Parallel()

	token1, err := GenerateRandomToken()
	if err != nil {
		t.Fatalf("GenerateRandomToken returned error: %v", err)
	}
	if token1 == "" {
		t.Fatalf("GenerateRandomToken returned empty token")
	}
	// 32 bytes -> 64 hex characters
	if len(token1) != 64 {
		t.Fatalf("expected token length 64, got %d", len(token1))
	}
	match, _ := regexp.MatchString("^[0-9a-f]+$", token1)
	if !match {
		t.Fatalf("token contains non-hex characters: %s", token1)
	}

	token2, err := GenerateRandomToken()
	if err != nil {
		t.Fatalf("GenerateRandomToken returned error on second call: %v", err)
	}
	if token1 == token2 {
		t.Fatalf("expected two generated tokens to differ")
	}
}

func TestValidateExpiry(t *testing.T) {
	t.Parallel()

	future := time.Now().Add(1 * time.Minute)
	if !ValidateExpiry(future) {
		t.Fatalf("ValidateExpiry returned false for future time")
	}

	past := time.Now().Add(-1 * time.Minute)
	if ValidateExpiry(past) {
		t.Fatalf("ValidateExpiry returned true for past time")
	}
}
