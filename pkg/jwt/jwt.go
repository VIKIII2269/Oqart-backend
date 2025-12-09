package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/oqart/backend/config"
)

type TokenManager struct {
	config *config.JWTConfig
}

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

func NewTokenManager(cfg *config.JWTConfig) *TokenManager {
	return &TokenManager{
		config: cfg,
	}
}

// GenerateTokenPair generates access and refresh tokens
func (tm *TokenManager) GenerateTokenPair(userID, email, role string) (*TokenPair, error) {
	// Generate access token
	accessToken, expiresAt, err := tm.GenerateAccessToken(userID, email, role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token
	refreshToken, err := tm.GenerateRefreshToken(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

// GenerateAccessToken generates a new access token
func (tm *TokenManager) GenerateAccessToken(userID, email, role string) (string, time.Time, error) {
	expiresAt := time.Now().Add(tm.config.Expiry)

	claims := &Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "oqart-backend",
			Subject:   userID,
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(tm.config.Secret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

// GenerateRefreshToken generates a new refresh token
func (tm *TokenManager) GenerateRefreshToken(userID string) (string, error) {
	expiresAt := time.Now().Add(tm.config.RefreshTokenExpiry)

	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "oqart-backend",
			Subject:   userID,
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(tm.config.Secret))
}

// ValidateToken validates and parses a token
func (tm *TokenManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tm.config.Secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// ExtractClaims extracts claims from a token without validation
func (tm *TokenManager) ExtractClaims(tokenString string) (*Claims, error) {
	token, _, err := jwt.NewParser().ParseUnverified(tokenString, &Claims{})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// RefreshAccessToken generates a new access token from a refresh token
func (tm *TokenManager) RefreshAccessToken(refreshToken string) (string, time.Time, error) {
	claims, err := tm.ValidateToken(refreshToken)
	if err != nil {
		return "", time.Time{}, err
	}

	// Generate new access token
	return tm.GenerateAccessToken(claims.UserID, claims.Email, claims.Role)
}

// GetTokenExpiry returns the expiry time of a token
func (tm *TokenManager) GetTokenExpiry(tokenString string) (time.Time, error) {
	claims, err := tm.ExtractClaims(tokenString)
	if err != nil {
		return time.Time{}, err
	}

	return claims.ExpiresAt.Time, nil
}

// IsTokenExpired checks if a token is expired
func (tm *TokenManager) IsTokenExpired(tokenString string) bool {
	expiresAt, err := tm.GetTokenExpiry(tokenString)
	if err != nil {
		return true
	}

	return time.Now().After(expiresAt)
}

// Custom errors
var (
	ErrInvalidToken = errors.New("invalid token")
	ErrTokenExpired = errors.New("token expired")
)
