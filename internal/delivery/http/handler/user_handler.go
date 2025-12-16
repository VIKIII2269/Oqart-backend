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
// @Description Get the profile of the currently authenticated user
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.UserResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /users/me [get]
func (h *UserHandler) GetMe(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": apperrors.NewErrorResponse(apperrors.ErrUnauthorized),
		})
		return
	}

	user, err := h.userUseCase.GetProfile(c.Request.Context(), userID)
	if err != nil {
		statusCode := apperrors.GetStatusCode(err)
		c.JSON(statusCode, gin.H{
			"error": apperrors.NewErrorResponse(err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    user,
	})
}

// UpdateMe godoc
// @Summary Update current user profile
// @Description Update the profile of the currently authenticated user
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpdateProfileRequest true "Profile update details"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
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

	user, err := h.userUseCase.UpdateProfile(c.Request.Context(), userID, &req)
	if err != nil {
		statusCode := apperrors.GetStatusCode(err)
		c.JSON(statusCode, gin.H{
			"error": apperrors.NewErrorResponse(err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    user,
		"message": "Profile updated successfully",
	})
}

// ChangePassword godoc
// @Summary Change password
// @Description Change the password for the currently authenticated user
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.ChangePasswordRequest true "Password change details"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
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

	response, err := h.userUseCase.ChangePassword(c.Request.Context(), userID, &req)
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

// UpdateEmail godoc
// @Summary Update email address
// @Description Update the email address for the currently authenticated user
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpdateEmailRequest true "Email update details"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 409 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
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

	response, err := h.userUseCase.UpdateEmail(c.Request.Context(), userID, &req)
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

// UpdatePhone godoc
// @Summary Update phone number
// @Description Update the phone number for the currently authenticated user
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpdatePhoneRequest true "Phone update details"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 409 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
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

	response, err := h.userUseCase.UpdatePhone(c.Request.Context(), userID, &req)
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

// DeleteAccount godoc
// @Summary Delete account
// @Description Delete the account of the currently authenticated user
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.MessageResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /users/me [delete]
func (h *UserHandler) DeleteAccount(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": apperrors.NewErrorResponse(apperrors.ErrUnauthorized),
		})
		return
	}

	response, err := h.userUseCase.DeleteAccount(c.Request.Context(), userID)
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

// GetPreferences godoc
// @Summary Get user preferences
// @Description Get the preferences for the currently authenticated user
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.UserPreferencesResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
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
// @Description Update the preferences for the currently authenticated user
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpdatePreferencesRequest true "Preferences update details"
// @Success 200 {object} dto.UserPreferencesResponse
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
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

// Address Management Handlers

// CreateAddress godoc
// @Summary Create a new address
// @Description Create a new address for the currently authenticated user
// @Tags Addresses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.AddressRequest true "Address details"
// @Success 201 {object} dto.AddressResponse
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /addresses [post]
func (h *UserHandler) CreateAddress(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": apperrors.NewErrorResponse(apperrors.ErrUnauthorized),
		})
		return
	}

	var req dto.AddressRequest
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

// GetAddresses godoc
// @Summary Get all addresses
// @Description Get all addresses for the currently authenticated user
// @Tags Addresses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.AddressResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /addresses [get]
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

// GetAddress godoc
// @Summary Get an address
// @Description Get a specific address by ID
// @Tags Addresses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Address ID"
// @Success 200 {object} dto.AddressResponse
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 403 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /addresses/{id} [get]
func (h *UserHandler) GetAddress(c *gin.Context) {
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

	address, err := h.userUseCase.GetAddress(c.Request.Context(), userID, addressID)
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
	})
}

// UpdateAddress godoc
// @Summary Update an address
// @Description Update a specific address by ID
// @Tags Addresses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Address ID"
// @Param request body dto.AddressRequest true "Address update details"
// @Success 200 {object} dto.AddressResponse
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 403 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /addresses/{id} [put]
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

	var req dto.AddressRequest
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
// @Summary Delete an address
// @Description Delete a specific address by ID
// @Tags Addresses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Address ID"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 403 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /addresses/{id} [delete]
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

	response, err := h.userUseCase.DeleteAddress(c.Request.Context(), userID, addressID)
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
