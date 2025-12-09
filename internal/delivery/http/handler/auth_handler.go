package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oqart/backend/internal/delivery/http/dto"
	"github.com/oqart/backend/internal/delivery/http/middleware"
	"github.com/oqart/backend/internal/usecase"
	apperrors "github.com/oqart/backend/pkg/errors"
	"github.com/oqart/backend/pkg/logger"
	"github.com/oqart/backend/pkg/validator"
)

type AuthHandler struct {
	authUseCase usecase.AuthUseCase
	logger      *logger.Logger
}

func NewAuthHandler(authUseCase usecase.AuthUseCase, logger *logger.Logger) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
		logger:      logger,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Registration details"
// @Success 201 {object} dto.AuthResponse
// @Failure 400 {object} apperrors.ErrorResponse
// @Failure 409 {object} apperrors.ErrorResponse
// @Failure 500 {object} apperrors.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": apperrors.NewErrorResponse(apperrors.ErrInvalidInput),
		})
		return
	}

	// Validate request
	if err := validator.Validate(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": apperrors.NewErrorResponseWithDetails(apperrors.ErrValidation, err.Error()),
		})
		return
	}

	// Register user
	response, err := h.authUseCase.Register(c.Request.Context(), &req)
	if err != nil {
		statusCode := apperrors.GetStatusCode(err)
		c.JSON(statusCode, gin.H{
			"error": apperrors.NewErrorResponse(err),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    response,
		"message": "User registered successfully",
	})
}

// Login godoc
// @Summary Login with email and password
// @Description Authenticate user with email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login credentials"
// @Success 200 {object} dto.AuthResponse
// @Failure 400 {object} apperrors.ErrorResponse
// @Failure 401 {object} apperrors.ErrorResponse
// @Failure 500 {object} apperrors.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": apperrors.NewErrorResponse(apperrors.ErrInvalidInput),
		})
		return
	}

	// Validate request
	if err := validator.Validate(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": apperrors.NewErrorResponseWithDetails(apperrors.ErrValidation, err.Error()),
		})
		return
	}

	// Login
	response, err := h.authUseCase.Login(c.Request.Context(), &req)
	if err != nil {
		statusCode := apperrors.GetStatusCode(err)
		c.JSON(statusCode, gin.H{
			"error": apperrors.NewErrorResponse(err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
		"message": "Login successful",
	})
}

// SendPhoneOTP godoc
// @Summary Send OTP to phone number
// @Description Send OTP for phone number authentication
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.PhoneLoginRequest true "Phone number"
// @Success 200 {object} dto.OTPResponse
// @Failure 400 {object} apperrors.ErrorResponse
// @Failure 429 {object} apperrors.ErrorResponse
// @Failure 500 {object} apperrors.ErrorResponse
// @Router /auth/login/phone [post]
func (h *AuthHandler) SendPhoneOTP(c *gin.Context) {
	var req dto.PhoneLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": apperrors.NewErrorResponse(apperrors.ErrInvalidInput),
		})
		return
	}

	// Validate request
	if err := validator.Validate(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": apperrors.NewErrorResponseWithDetails(apperrors.ErrValidation, err.Error()),
		})
		return
	}

	// Send OTP
	response, err := h.authUseCase.SendPhoneOTP(c.Request.Context(), &req)
	if err != nil {
		statusCode := apperrors.GetStatusCode(err)
		c.JSON(statusCode, gin.H{
			"error": apperrors.NewErrorResponse(err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// VerifyPhoneOTP godoc
// @Summary Verify OTP and login
// @Description Verify OTP code and authenticate user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.VerifyOTPRequest true "Phone and OTP code"
// @Success 200 {object} dto.AuthResponse
// @Failure 400 {object} apperrors.ErrorResponse
// @Failure 500 {object} apperrors.ErrorResponse
// @Router /auth/verify-otp [post]
func (h *AuthHandler) VerifyPhoneOTP(c *gin.Context) {
	var req dto.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": apperrors.NewErrorResponse(apperrors.ErrInvalidInput),
		})
		return
	}

	// Validate request
	if err := validator.Validate(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": apperrors.NewErrorResponseWithDetails(apperrors.ErrValidation, err.Error()),
		})
		return
	}

	// Verify OTP
	response, err := h.authUseCase.VerifyPhoneOTP(c.Request.Context(), &req)
	if err != nil {
		statusCode := apperrors.GetStatusCode(err)
		c.JSON(statusCode, gin.H{
			"error": apperrors.NewErrorResponse(err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
		"message": "OTP verified successfully",
	})
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Get new access token using refresh token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} dto.AuthResponse
// @Failure 400 {object} apperrors.ErrorResponse
// @Failure 401 {object} apperrors.ErrorResponse
// @Failure 500 {object} apperrors.ErrorResponse
// @Router /auth/refresh-token [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": apperrors.NewErrorResponse(apperrors.ErrInvalidInput),
		})
		return
	}

	// Validate request
	if err := validator.Validate(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": apperrors.NewErrorResponseWithDetails(apperrors.ErrValidation, err.Error()),
		})
		return
	}

	// Refresh token
	response, err := h.authUseCase.RefreshToken(c.Request.Context(), &req)
	if err != nil {
		statusCode := apperrors.GetStatusCode(err)
		c.JSON(statusCode, gin.H{
			"error": apperrors.NewErrorResponse(err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
		"message": "Token refreshed successfully",
	})
}

// ForgotPassword godoc
// @Summary Request password reset
// @Description Send password reset link to email
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.ForgotPasswordRequest true "Email address"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} apperrors.ErrorResponse
// @Failure 500 {object} apperrors.ErrorResponse
// @Router /auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req dto.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": apperrors.NewErrorResponse(apperrors.ErrInvalidInput),
		})
		return
	}

	// Validate request
	if err := validator.Validate(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": apperrors.NewErrorResponseWithDetails(apperrors.ErrValidation, err.Error()),
		})
		return
	}

	// Request password reset
	response, err := h.authUseCase.ForgotPassword(c.Request.Context(), &req)
	if err != nil {
		statusCode := apperrors.GetStatusCode(err)
		c.JSON(statusCode, gin.H{
			"error": apperrors.NewErrorResponse(err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// ResetPassword godoc
// @Summary Reset password
// @Description Reset password using reset token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.ResetPasswordRequest true "Reset token and new password"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} apperrors.ErrorResponse
// @Failure 500 {object} apperrors.ErrorResponse
// @Router /auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": apperrors.NewErrorResponse(apperrors.ErrInvalidInput),
		})
		return
	}

	// Validate request
	if err := validator.Validate(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": apperrors.NewErrorResponseWithDetails(apperrors.ErrValidation, err.Error()),
		})
		return
	}

	// Reset password
	response, err := h.authUseCase.ResetPassword(c.Request.Context(), &req)
	if err != nil {
		statusCode := apperrors.GetStatusCode(err)
		c.JSON(statusCode, gin.H{
			"error": apperrors.NewErrorResponse(err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// VerifyEmail godoc
// @Summary Verify email address
// @Description Verify email using verification token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param token query string true "Verification token"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} apperrors.ErrorResponse
// @Failure 500 {object} apperrors.ErrorResponse
// @Router /auth/verify-email [get]
func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": apperrors.NewErrorResponse(apperrors.ErrInvalidInput),
		})
		return
	}

	// Verify email
	response, err := h.authUseCase.VerifyEmail(c.Request.Context(), token)
	if err != nil {
		statusCode := apperrors.GetStatusCode(err)
		c.JSON(statusCode, gin.H{
			"error": apperrors.NewErrorResponse(err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// Logout godoc
// @Summary Logout user
// @Description Logout user and invalidate all sessions
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.MessageResponse
// @Failure 401 {object} apperrors.ErrorResponse
// @Failure 500 {object} apperrors.ErrorResponse
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": apperrors.NewErrorResponse(apperrors.ErrUnauthorized),
		})
		return
	}

	// Logout
	if err := h.authUseCase.Logout(c.Request.Context(), userID); err != nil {
		statusCode := apperrors.GetStatusCode(err)
		c.JSON(statusCode, gin.H{
			"error": apperrors.NewErrorResponse(err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Logged out successfully",
	})
}
