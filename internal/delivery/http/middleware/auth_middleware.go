package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/oqart/backend/internal/domain"
	"github.com/oqart/backend/internal/repository/postgres"
	apperrors "github.com/oqart/backend/pkg/errors"
	"github.com/oqart/backend/pkg/jwt"
)

const (
	AuthorizationHeader = "Authorization"
	UserContextKey      = "user"
	UserIDContextKey    = "user_id"
	UserRoleContextKey  = "user_role"
)

type AuthMiddleware struct {
	jwtManager *jwt.TokenManager
	userRepo   postgres.UserRepository
}

func NewAuthMiddleware(jwtManager *jwt.TokenManager, userRepo postgres.UserRepository) *AuthMiddleware {
	return &AuthMiddleware{
		jwtManager: jwtManager,
		userRepo:   userRepo,
	}
}

// RequireAuth middleware validates JWT token
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := extractToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": apperrors.NewErrorResponse(apperrors.ErrUnauthorized),
			})
			c.Abort()
			return
		}

		// Validate token
		claims, err := m.jwtManager.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": apperrors.NewErrorResponse(err),
			})
			c.Abort()
			return
		}

		// Parse user ID
		userID, err := uuid.Parse(claims.UserID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": apperrors.NewErrorResponse(apperrors.ErrInvalidToken),
			})
			c.Abort()
			return
		}

		// Get user from database
		user, err := m.userRepo.GetByID(c.Request.Context(), userID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": apperrors.NewErrorResponse(apperrors.ErrUnauthorized),
			})
			c.Abort()
			return
		}

		// Check if user is active
		if !user.IsActive() {
			c.JSON(http.StatusForbidden, gin.H{
				"error": apperrors.NewErrorResponse(apperrors.New("ACCOUNT_INACTIVE", "Your account is inactive", http.StatusForbidden)),
			})
			c.Abort()
			return
		}

		// Set user in context
		c.Set(UserContextKey, user)
		c.Set(UserIDContextKey, userID)
		c.Set(UserRoleContextKey, user.Role)

		c.Next()
	}
}

// RequireRole middleware checks if user has required role
func (m *AuthMiddleware) RequireRole(roles ...domain.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get(UserContextKey)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": apperrors.NewErrorResponse(apperrors.ErrUnauthorized),
			})
			c.Abort()
			return
		}

		u := user.(*domain.User)

		// Check if user has any of the required roles
		hasRole := false
		for _, role := range roles {
			if u.Role == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{
				"error": apperrors.NewErrorResponse(apperrors.ErrPermissionDenied),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAdmin middleware ensures user is admin or super admin
func (m *AuthMiddleware) RequireAdmin() gin.HandlerFunc {
	return m.RequireRole(domain.RoleAdmin, domain.RoleSuperAdmin)
}

// RequireVendor middleware ensures user is a vendor
func (m *AuthMiddleware) RequireVendor() gin.HandlerFunc {
	return m.RequireRole(domain.RoleVendor)
}

// OptionalAuth middleware tries to authenticate but doesn't fail if no token
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := extractToken(c)
		if err != nil {
			c.Next()
			return
		}

		// Validate token
		claims, err := m.jwtManager.ValidateToken(token)
		if err != nil {
			c.Next()
			return
		}

		// Parse user ID
		userID, err := uuid.Parse(claims.UserID)
		if err != nil {
			c.Next()
			return
		}

		// Get user from database
		user, err := m.userRepo.GetByID(c.Request.Context(), userID)
		if err != nil {
			c.Next()
			return
		}

		// Set user in context if found
		c.Set(UserContextKey, user)
		c.Set(UserIDContextKey, userID)
		c.Set(UserRoleContextKey, user.Role)

		c.Next()
	}
}

// extractToken extracts JWT token from Authorization header
func extractToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader(AuthorizationHeader)
	if authHeader == "" {
		return "", apperrors.ErrUnauthorized
	}

	// Check if header starts with "Bearer "
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", apperrors.ErrUnauthorized
	}

	return parts[1], nil
}

// Helper functions to get user from context

func GetUser(c *gin.Context) (*domain.User, error) {
	user, exists := c.Get(UserContextKey)
	if !exists {
		return nil, apperrors.ErrUnauthorized
	}
	return user.(*domain.User), nil
}

func GetUserID(c *gin.Context) (uuid.UUID, error) {
	userID, exists := c.Get(UserIDContextKey)
	if !exists {
		return uuid.Nil, apperrors.ErrUnauthorized
	}
	return userID.(uuid.UUID), nil
}

func GetUserRole(c *gin.Context) (domain.UserRole, error) {
	role, exists := c.Get(UserRoleContextKey)
	if !exists {
		return "", apperrors.ErrUnauthorized
	}
	return role.(domain.UserRole), nil
}
