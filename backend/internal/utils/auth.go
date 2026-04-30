package utils

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a plain text password using bcrypt
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckPasswordHash compares a password with a hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateSessionID generates a new UUID for session
func GenerateSessionID() (string, error) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return uuid.String(), nil
}

// GenerateRandomToken generates a secure random token (alternative to UUID)
func GenerateRandomToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", bytes), nil
}

// ValidateExpiry checks if a time is still valid
func ValidateExpiry(expiresAt time.Time) bool {
	return time.Now().Before(expiresAt)
}
