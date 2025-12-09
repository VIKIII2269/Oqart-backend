package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const (
	DefaultBcryptCost = 12
)

// HashPassword hashes a password using bcrypt
func HashPassword(password string, cost int) (string, error) {
	if cost == 0 {
		cost = DefaultBcryptCost
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(bytes), nil
}

// CheckPassword checks if a password matches its hash
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateRandomToken generates a random token
func GenerateRandomToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random token: %w", err)
	}

	return base64.URLEncoding.EncodeToString(bytes), nil
}

// GenerateSecureToken generates a secure random token (32 bytes)
func GenerateSecureToken() (string, error) {
	return GenerateRandomToken(32)
}
