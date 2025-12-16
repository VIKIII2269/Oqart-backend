package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/oqart/backend/internal/delivery/http/dto"
	"github.com/oqart/backend/internal/delivery/http/middleware"
	"github.com/oqart/backend/internal/usecase"
	apperrors "github.com/oqart/backend/pkg/errors"
	"github.com/oqart/backend/pkg/logger"
	"github.com/oqart/backend/pkg/validator"
)

type UserHandler struct {
	userUseCase usecase.UserUseCase
	logger      *logger.Logger
}

func NewUserHandler(userUseCase usecase.UserUseCase, logger *logger.Logger) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
		logger:      logger,
	}
}

// GetMe godoc
// @Summary Get current user profile
// @Description Get the authenticated user's profile information
// @Tags User Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.UserResponse
// @Failure 401 {object} apperrors.ErrorResponse
// @Failure 500 {object} apperrors.ErrorResponse
// @Router /users/me [get]
func (h *UserHandler) GetMe(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": apperrors.NewErrorResponse(apperrors.ErrUnauthorized),
		})
		return
	}

	profile, err := h.userUseCase.GetProfile(c.Request.Context(), userID)
	if err != nil {
		statusCode := apperrors.GetStatusCode(err)
		c.JSON(statusCode, gin.H{
			"error": apperrors.NewErrorResponse(err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    profile,
	})
}

// UpdateMe godoc
// @Summary Update current user profile
// @Description Update the authenticated user's profile information
// @Tags User Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpdateProfileRequest true "Profile update data"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} apperrors.ErrorResponse
// @Failure 401 {object} apperrors.ErrorResponse
// @Failure 500 {object} apperrors.ErrorResponse
// @Router /users/me [put]
func (h *UserHandler) UpdateMe(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": apperrors.NewErrorResponse(apperrors.ErrUnauthorized),
		})
		return
	}

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": apperrors.NewErrorResponse(apperrors.ErrInvalidInput),
		})
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": apperrors.NewErrorResponseWithDetails(apperrors.ErrValidation, err.Error()),
		})
		return
	}

	profile, err := h.userUseCase.UpdateProfile(c.Request.Context(), userID, &req)
	if err != nil {
		statusCode := apperrors.GetStatusCode(err)
		c.JSON(statusCode, gin.H{
			"error": apperrors.NewErrorResponse(err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    profile,
		"message": "Profile updated successfully",
	})
}

// ChangePassword godoc
// @Summary Change user password
// @Description Change the authenticated user's password
// @Tags User Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.ChangePasswordRequest true "Password change data"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} apperrors.ErrorResponse
// @Failure 401 {object} apperrors.ErrorResponse
// @Failure 500 {object} apperrors.ErrorResponse
// @Router /users/me/password [put]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": apperrors.NewErrorResponse(apperrors.ErrUnauthorized),
		})
		return
	}

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": apperrors.NewErrorResponse(apperrors.ErrInvalidInput),
		})
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": apperrors.NewErrorResponseWithDetails(apperrors.ErrValidation, err.Error()),
		})
		return
	}

	if err := h.userUseCase.ChangePassword(c.Request.Context(), userID, &req); err != nil {
		statusCode := apperrors.GetStatusCode(err)
		c.JSON(statusCode, gin.H{
			"error": apperrors.NewErrorResponse(err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Password changed successfully",
	})
}

// UpdateEmail godoc
// @Summary Update user email
// @Description Update the authenticated user's email address
// @Tags User Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpdateEmailRequest true "Email update data"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} apperrors.ErrorResponse
// @Failure 401 {object} apperrors.ErrorResponse
// @Failure 409 {object} apperrors.ErrorResponse
// @Failure 500 {object} apperrors.ErrorResponse
// @Router /users/me/email [put]
func (h *UserHandler) UpdateEmail(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": apperrors.NewErrorResponse(apperrors.ErrUnauthorized),
		})
		return
	}

	var req dto.UpdateEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": apperrors.NewErrorResponse(apperrors.ErrInvalidInput),
		})
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": apperrors.NewErrorResponseWithDetails(apperrors.ErrValidation, err.Error()),
		})
		return
	}

	if err := h.userUseCase.UpdateEmail(c.Request.Context(), userID, &req); err != nil {
		statusCode := apperrors.GetStatusCode(err)
		c.JSON(statusCode, gin.H{
			"error": apperrors.NewErrorResponse(err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Email updated successfully. Please verify your new email.",
	})
}

// UpdatePhone godoc
// @Summary Update user phone number
// @Description Update the authenticated user's phone number
// @Tags User Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpdatePhoneRequest true "Phone update data"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} apperrors.ErrorResponse
// @Failure 401 {object} apperrors.ErrorResponse
// @Failure 409 {object} apperrors.ErrorResponse
// @Failure 500 {object} apperrors.ErrorResponse
// @Router /users/me/phone [put]
func (h *UserHandler) UpdatePhone(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": apperrors.NewErrorResponse(apperrors.ErrUnauthorized),
		})
		return
	}

	var req dto.UpdatePhoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": apperrors.NewErrorResponse(apperrors.ErrInvalidInput),
		})
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": apperrors.NewErrorResponseWithDetails(apperrors.ErrValidation, err.Error()),
		})
		return
	}

	if err := h.userUseCase.UpdatePhone(c.Request.Context(), userID, &req); err != nil {
		statusCode := apperrors.GetStatusCode(err)
		c.JSON(statusCode, gin.H{
			"error": apperrors.NewErrorResponse(err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Phone number updated successfully. Please verify your new phone number.",
	})
}

// DeleteAccount godoc
// @Summary Delete user account
// @Description Soft delete the authenticated user's account
// @Tags User Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.MessageResponse
// @Failure 401 {object} apperrors.ErrorResponse
// @Failure 500 {object} apperrors.ErrorResponse
// @Router /users/me [delete]
func (h *UserHandler) DeleteAccount(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": apperrors.NewErrorResponse(apperrors.ErrUnauthorized),
		})
		return
	}

	if err := h.userUseCase.DeleteAccount(c.Request.Context(), userID); err != nil {
		statusCode := apperrors.GetStatusCode(err)
		c.JSON(statusCode, gin.H{
			"error": apperrors.NewErrorResponse(err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Account deleted successfully",
	})
}

// GetPreferences godoc
// @Summary Get user preferences
// @Description Get the authenticated user's preferences
// @Tags User Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.PreferencesResponse
// @Failure 401 {object} apperrors.ErrorResponse
// @Failure 500 {object} apperrors.ErrorResponse
// @Router /users/me/preferences [get]
func (h *UserHandler) GetPreferences(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": apperrors.NewErrorResponse(apperrors.ErrUnauthorized),
		})
		return
	}

	prefs, err := h.userUseCase.GetPreferences(c.Request.Context(), userID)
	if err != nil {
		statusCode := apperrors.GetStatusCode(err)
		c.JSON(statusCode, gin.H{
			"error": apperrors.NewErrorResponse(err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    prefs,
	})
}

// UpdatePreferences godoc
// @Summary Update user preferences
// @Description Update the authenticated user's preferences
// @Tags User Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpdatePreferencesRequest true "Preferences update data"
// @Success 200 {object} dto.PreferencesResponse
// @Failure 400 {object} apperrors.ErrorResponse
// @Failure 401 {object} apperrors.ErrorResponse
// @Failure 500 {object} apperrors.ErrorResponse
// @Router /users/me/preferences [put]
func (h *UserHandler) UpdatePreferences(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": apperrors.NewErrorResponse(apperrors.ErrUnauthorized),
		})
		return
	}

	var req dto.UpdatePreferencesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": apperrors.NewErrorResponse(apperrors.ErrInvalidInput),
		})
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": apperrors.NewErrorResponseWithDetails(apperrors.ErrValidation, err.Error()),
		})
		return
	}

	prefs, err := h.userUseCase.UpdatePreferences(c.Request.Context(), userID, &req)
	if err != nil {
		statusCode := apperrors.GetStatusCode(err)
		c.JSON(statusCode, gin.H{
			"error": apperrors.NewErrorResponse(err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    prefs,
		"message": "Preferences updated successfully",
	})
}

// GetAddresses godoc
// @Summary Get user addresses
// @Description Get all addresses for the authenticated user
// @Tags Address Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.AddressResponse
// @Failure 401 {object} apperrors.ErrorResponse
// @Failure 500 {object} apperrors.ErrorResponse
// @Router /users/me/addresses [get]
func (h *UserHandler) GetAddresses(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": apperrors.NewErrorResponse(apperrors.ErrUnauthorized),
		})
		return
	}

	addresses, err := h.userUseCase.GetAddresses(c.Request.Context(), userID)
	if err != nil {
		statusCode := apperrors.GetStatusCode(err)
		c.JSON(statusCode, gin.H{
			"error": apperrors.NewErrorResponse(err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    addresses,
	})
}

// CreateAddress godoc
// @Summary Create new address
// @Description Create a new address for the authenticated user
// @Tags Address Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateAddressRequest true "Address data"
// @Success 201 {object} dto.AddressResponse
// @Failure 400 {object} apperrors.ErrorResponse
// @Failure 401 {object} apperrors.ErrorResponse
// @Failure 500 {object} apperrors.ErrorResponse
// @Router /users/me/addresses [post]
func (h *UserHandler) CreateAddress(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": apperrors.NewErrorResponse(apperrors.ErrUnauthorized),
		})
		return
	}

	var req dto.CreateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": apperrors.NewErrorResponse(apperrors.ErrInvalidInput),
		})
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": apperrors.NewErrorResponseWithDetails(apperrors.ErrValidation, err.Error()),
		})
		return
	}

	address, err := h.userUseCase.CreateAddress(c.Request.Context(), userID, &req)
	if err != nil {
		statusCode := apperrors.GetStatusCode(err)
		c.JSON(statusCode, gin.H{
			"error": apperrors.NewErrorResponse(err),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    address,
		"message": "Address created successfully",
	})
}

// UpdateAddress godoc
// @Summary Update address
// @Description Update an existing address for the authenticated user
// @Tags Address Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Address ID"
// @Param request body dto.UpdateAddressRequest true "Address update data"
// @Success 200 {object} dto.AddressResponse
// @Failure 400 {object} apperrors.ErrorResponse
// @Failure 401 {object} apperrors.ErrorResponse
// @Failure 403 {object} apperrors.ErrorResponse
// @Failure 404 {object} apperrors.ErrorResponse
// @Failure 500 {object} apperrors.ErrorResponse
// @Router /users/me/addresses/{id} [put]
func (h *UserHandler) UpdateAddress(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": apperrors.NewErrorResponse(apperrors.ErrUnauthorized),
		})
		return
	}

	addressID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": apperrors.NewErrorResponse(apperrors.ErrInvalidInput),
		})
		return
	}

	var req dto.UpdateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": apperrors.NewErrorResponse(apperrors.ErrInvalidInput),
		})
		return
	}

	if err := validator.Validate(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": apperrors.NewErrorResponseWithDetails(apperrors.ErrValidation, err.Error()),
		})
		return
	}

	address, err := h.userUseCase.UpdateAddress(c.Request.Context(), userID, addressID, &req)
	if err != nil {
		statusCode := apperrors.GetStatusCode(err)
		c.JSON(statusCode, gin.H{
			"error": apperrors.NewErrorResponse(err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    address,
		"message": "Address updated successfully",
	})
}

// DeleteAddress godoc
// @Summary Delete address
// @Description Delete an address for the authenticated user
// @Tags Address Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Address ID"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} apperrors.ErrorResponse
// @Failure 401 {object} apperrors.ErrorResponse
// @Failure 403 {object} apperrors.ErrorResponse
// @Failure 404 {object} apperrors.ErrorResponse
// @Failure 500 {object} apperrors.ErrorResponse
// @Router /users/me/addresses/{id} [delete]
func (h *UserHandler) DeleteAddress(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": apperrors.NewErrorResponse(apperrors.ErrUnauthorized),
		})
		return
	}

	addressID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": apperrors.NewErrorResponse(apperrors.ErrInvalidInput),
		})
		return
	}

	if err := h.userUseCase.DeleteAddress(c.Request.Context(), userID, addressID); err != nil {
		statusCode := apperrors.GetStatusCode(err)
		c.JSON(statusCode, gin.H{
			"error": apperrors.NewErrorResponse(err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Address deleted successfully",
	})
}

// SetDefaultAddress godoc
// @Summary Set default address
// @Description Set an address as the default for the authenticated user
// @Tags Address Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Address ID"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} apperrors.ErrorResponse
// @Failure 401 {object} apperrors.ErrorResponse
// @Failure 403 {object} apperrors.ErrorResponse
// @Failure 404 {object} apperrors.ErrorResponse
// @Failure 500 {object} apperrors.ErrorResponse
// @Router /users/me/addresses/{id}/default [put]
func (h *UserHandler) SetDefaultAddress(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": apperrors.NewErrorResponse(apperrors.ErrUnauthorized),
		})
		return
	}

	addressID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": apperrors.NewErrorResponse(apperrors.ErrInvalidInput),
		})
		return
	}

	if err := h.userUseCase.SetDefaultAddress(c.Request.Context(), userID, addressID); err != nil {
		statusCode := apperrors.GetStatusCode(err)
		c.JSON(statusCode, gin.H{
			"error": apperrors.NewErrorResponse(err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Default address set successfully",
	})
}
