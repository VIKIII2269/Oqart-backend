package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// GenerateOTP generates a random OTP of specified length
func GenerateOTP(length int) (string, error) {
	if length <= 0 {
		length = 6
	}

	// Calculate max value (10^length - 1)
	max := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(length)), nil)

	// Generate random number
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", fmt.Errorf("failed to generate OTP: %w", err)
	}

	// Format with leading zeros
	format := fmt.Sprintf("%%0%dd", length)
	return fmt.Sprintf(format, n), nil
}

// Generate6DigitOTP generates a 6-digit OTP
func Generate6DigitOTP() (string, error) {
	return GenerateOTP(6)
}

// Generate4DigitOTP generates a 4-digit OTP
func Generate4DigitOTP() (string, error) {
	return GenerateOTP(4)
}
