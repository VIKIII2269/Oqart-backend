package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/oqart/backend/internal/delivery/http/dto"
	"github.com/oqart/backend/internal/delivery/http/middleware"
	"github.com/oqart/backend/internal/usecase"
	apperrors "github.com/oqart/backend/pkg/errors"
	"github.com/oqart/backend/pkg/logger"
	"github.com/oqart/backend/pkg/validator"
)

type VendorHandler struct {
	vendorUseCase usecase.VendorUseCase
	logger        *logger.Logger
}

func NewVendorHandler(vendorUseCase usecase.VendorUseCase, logger *logger.Logger) *VendorHandler {
	return &VendorHandler{
		vendorUseCase: vendorUseCase,
		logger:        logger,
	}
}

// OnboardStep1 godoc
// @Summary Vendor onboarding - Step 1 (Business Details)
// @Description Complete the first step of vendor onboarding by providing business details
// @Tags Vendor Onboarding
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.VendorOnboardingStep1Request true "Business details"
// @Success 201 {object} dto.VendorResponse
// @Failure 400 {object} apperrors.ErrorResponse "Invalid input or validation error"
// @Failure 401 {object} apperrors.ErrorResponse "Unauthorized"
// @Failure 409 {object} apperrors.ErrorResponse "Vendor already exists"
// @Failure 500 {object} apperrors.ErrorResponse "Internal server error"
// @Router /vendors/onboard/step1 [post]
func (h *VendorHandler) OnboardStep1(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": apperrors.NewErrorResponse(apperrors.ErrUnauthorized),
		})
		return
	}

	var req dto.VendorOnboardingStep1Request
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

	vendor, err := h.vendorUseCase.OnboardStep1(c.Request.Context(), userID, &req)
	if err != nil {
		statusCode := apperrors.GetStatusCode(err)
		c.JSON(statusCode, gin.H{
			"error": apperrors.NewErrorResponse(err),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    vendor,
		"message": "Step 1 completed successfully. Please proceed to step 2.",
	})
}

// OnboardStep2 godoc
// @Summary Vendor onboarding - Step 2 (Business Address)
// @Description Complete the second step of vendor onboarding by providing business address
// @Tags Vendor Onboarding
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.VendorOnboardingStep2Request true "Business address"
// @Success 201 {object} dto.VendorAddressResponse
// @Failure 400 {object} apperrors.ErrorResponse "Invalid input or validation error"
// @Failure 401 {object} apperrors.ErrorResponse "Unauthorized"
// @Failure 404 {object} apperrors.ErrorResponse "Vendor not found"
// @Failure 500 {object} apperrors.ErrorResponse "Internal server error"
// @Router /vendors/onboard/step2 [post]
func (h *VendorHandler) OnboardStep2(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": apperrors.NewErrorResponse(apperrors.ErrUnauthorized),
		})
		return
	}

	var req dto.VendorOnboardingStep2Request
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

	address, err := h.vendorUseCase.OnboardStep2(c.Request.Context(), userID, &req)
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
		"message": "Step 2 completed successfully. Please proceed to step 3.",
	})
}

// OnboardStep3 godoc
// @Summary Vendor onboarding - Step 3 (Bank Details)
// @Description Complete the third step of vendor onboarding by providing bank details
// @Tags Vendor Onboarding
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.VendorOnboardingStep3Request true "Bank details"
// @Success 201 {object} dto.VendorBankDetailsResponse
// @Failure 400 {object} apperrors.ErrorResponse "Invalid input or validation error"
// @Failure 401 {object} apperrors.ErrorResponse "Unauthorized"
// @Failure 404 {object} apperrors.ErrorResponse "Vendor not found"
// @Failure 500 {object} apperrors.ErrorResponse "Internal server error"
// @Router /vendors/onboard/step3 [post]
func (h *VendorHandler) OnboardStep3(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": apperrors.NewErrorResponse(apperrors.ErrUnauthorized),
		})
		return
	}

	var req dto.VendorOnboardingStep3Request
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

	bankDetails, err := h.vendorUseCase.OnboardStep3(c.Request.Context(), userID, &req)
	if err != nil {
		statusCode := apperrors.GetStatusCode(err)
		c.JSON(statusCode, gin.H{
			"error": apperrors.NewErrorResponse(err),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    bankDetails,
		"message": "Step 3 completed successfully. Please proceed to step 4.",
	})
}

// OnboardStep4 godoc
// @Summary Vendor onboarding - Step 4 (Contact Person)
// @Description Complete the final step of vendor onboarding by providing contact person details
// @Tags Vendor Onboarding
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.VendorOnboardingStep4Request true "Contact person details"
// @Success 201 {object} dto.VendorContactPersonResponse
// @Failure 400 {object} apperrors.ErrorResponse "Invalid input or validation error"
// @Failure 401 {object} apperrors.ErrorResponse "Unauthorized"
// @Failure 404 {object} apperrors.ErrorResponse "Vendor not found"
// @Failure 500 {object} apperrors.ErrorResponse "Internal server error"
// @Router /vendors/onboard/step4 [post]
func (h *VendorHandler) OnboardStep4(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": apperrors.NewErrorResponse(apperrors.ErrUnauthorized),
		})
		return
	}

	var req dto.VendorOnboardingStep4Request
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

	contact, err := h.vendorUseCase.OnboardStep4(c.Request.Context(), userID, &req)
	if err != nil {
		statusCode := apperrors.GetStatusCode(err)
		c.JSON(statusCode, gin.H{
			"error": apperrors.NewErrorResponse(err),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    contact,
		"message": "Onboarding completed successfully! Your vendor account is now under review.",
	})
}

// GetMyVendor godoc
// @Summary Get vendor dashboard
// @Description Get the authenticated vendor's complete profile with all details
// @Tags Vendor Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.VendorDetailedResponse
// @Failure 401 {object} apperrors.ErrorResponse "Unauthorized"
// @Failure 404 {object} apperrors.ErrorResponse "Vendor not found"
// @Failure 500 {object} apperrors.ErrorResponse "Internal server error"
// @Router /vendors/me [get]
func (h *VendorHandler) GetMyVendor(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": apperrors.NewErrorResponse(apperrors.ErrUnauthorized),
		})
		return
	}

	vendor, err := h.vendorUseCase.GetMyVendor(c.Request.Context(), userID)
	if err != nil {
		statusCode := apperrors.GetStatusCode(err)
		c.JSON(statusCode, gin.H{
			"error": apperrors.NewErrorResponse(err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    vendor,
	})
}

// UpdateProfile godoc
// @Summary Update vendor profile
// @Description Update the authenticated vendor's profile information
// @Tags Vendor Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpdateVendorProfileRequest true "Profile update data"
// @Success 200 {object} dto.VendorResponse
// @Failure 400 {object} apperrors.ErrorResponse "Invalid input or validation error"
// @Failure 401 {object} apperrors.ErrorResponse "Unauthorized"
// @Failure 404 {object} apperrors.ErrorResponse "Vendor not found"
// @Failure 500 {object} apperrors.ErrorResponse "Internal server error"
// @Router /vendors/me [put]
func (h *VendorHandler) UpdateProfile(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": apperrors.NewErrorResponse(apperrors.ErrUnauthorized),
		})
		return
	}

	var req dto.UpdateVendorProfileRequest
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

	vendor, err := h.vendorUseCase.UpdateProfile(c.Request.Context(), userID, &req)
	if err != nil {
		statusCode := apperrors.GetStatusCode(err)
		c.JSON(statusCode, gin.H{
			"error": apperrors.NewErrorResponse(err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    vendor,
		"message": "Profile updated successfully",
	})
}

// GetVendorByID godoc
// @Summary Get vendor by ID
// @Description Get public vendor profile by vendor ID (only approved vendors)
// @Tags Vendors
// @Accept json
// @Produce json
// @Param id path string true "Vendor ID (UUID)"
// @Success 200 {object} dto.VendorResponse
// @Failure 400 {object} apperrors.ErrorResponse "Invalid vendor ID"
// @Failure 404 {object} apperrors.ErrorResponse "Vendor not found"
// @Failure 500 {object} apperrors.ErrorResponse "Internal server error"
// @Router /vendors/{id} [get]
func (h *VendorHandler) GetVendorByID(c *gin.Context) {
	vendorID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": apperrors.NewErrorResponse(apperrors.ErrInvalidInput),
		})
		return
	}

	vendor, err := h.vendorUseCase.GetVendorByID(c.Request.Context(), vendorID)
	if err != nil {
		statusCode := apperrors.GetStatusCode(err)
		c.JSON(statusCode, gin.H{
			"error": apperrors.NewErrorResponse(err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    vendor,
	})
}

// ListVendors godoc
// @Summary List vendors
// @Description Get a paginated list of vendors with optional filters
// @Tags Vendors
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)" default(1)
// @Param page_size query int false "Page size (default: 20, max: 100)" default(20)
// @Param status query string false "Filter by status (pending, under_review, approved, rejected, suspended, inactive)"
// @Param is_verified query bool false "Filter by verification status"
// @Param is_featured query bool false "Filter by featured status"
// @Param search query string false "Search by business name or description"
// @Success 200 {object} dto.VendorListResponse
// @Failure 400 {object} apperrors.ErrorResponse "Invalid parameters"
// @Failure 500 {object} apperrors.ErrorResponse "Internal server error"
// @Router /vendors [get]
func (h *VendorHandler) ListVendors(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	filters := make(map[string]interface{})

	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if isVerified := c.Query("is_verified"); isVerified != "" {
		filters["is_verified"] = isVerified == "true"
	}
	if isFeatured := c.Query("is_featured"); isFeatured != "" {
		filters["is_featured"] = isFeatured == "true"
	}
	if search := c.Query("search"); search != "" {
		filters["search"] = search
	}

	vendors, err := h.vendorUseCase.ListVendors(c.Request.Context(), filters, page, pageSize)
	if err != nil {
		statusCode := apperrors.GetStatusCode(err)
		c.JSON(statusCode, gin.H{
			"error": apperrors.NewErrorResponse(err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    vendors,
	})
}

// GetDocuments godoc
// @Summary Get vendor documents
// @Description Get all documents uploaded by the authenticated vendor
// @Tags Vendor Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.VendorDocumentResponse
// @Failure 401 {object} apperrors.ErrorResponse "Unauthorized"
// @Failure 404 {object} apperrors.ErrorResponse "Vendor not found"
// @Failure 500 {object} apperrors.ErrorResponse "Internal server error"
// @Router /vendors/me/documents [get]
func (h *VendorHandler) GetDocuments(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": apperrors.NewErrorResponse(apperrors.ErrUnauthorized),
		})
		return
	}

	documents, err := h.vendorUseCase.GetDocuments(c.Request.Context(), userID)
	if err != nil {
		statusCode := apperrors.GetStatusCode(err)
		c.JSON(statusCode, gin.H{
			"error": apperrors.NewErrorResponse(err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    documents,
	})
}

// ApproveVendor godoc
// @Summary Approve or reject vendor (Admin only)
// @Description Admin endpoint to approve, reject, or suspend a vendor
// @Tags Admin - Vendor Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Vendor ID (UUID)"
// @Param request body dto.AdminVendorApprovalRequest true "Approval decision"
// @Success 200 {object} map[string]interface{} "Vendor status updated successfully"
// @Failure 400 {object} apperrors.ErrorResponse "Invalid input or validation error"
// @Failure 401 {object} apperrors.ErrorResponse "Unauthorized"
// @Failure 403 {object} apperrors.ErrorResponse "Forbidden - Admin access required"
// @Failure 404 {object} apperrors.ErrorResponse "Vendor not found"
// @Failure 500 {object} apperrors.ErrorResponse "Internal server error"
// @Router /admin/vendors/{id}/approve [post]
func (h *VendorHandler) ApproveVendor(c *gin.Context) {
	adminID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": apperrors.NewErrorResponse(apperrors.ErrUnauthorized),
		})
		return
	}

	vendorID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": apperrors.NewErrorResponse(apperrors.ErrInvalidInput),
		})
		return
	}

	var req dto.AdminVendorApprovalRequest
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

	if err := h.vendorUseCase.ApproveVendor(c.Request.Context(), vendorID, adminID, &req); err != nil {
		statusCode := apperrors.GetStatusCode(err)
		c.JSON(statusCode, gin.H{
			"error": apperrors.NewErrorResponse(err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Vendor status updated successfully",
	})
}

// VerifyDocument godoc
// @Summary Verify or reject vendor document (Admin only)
// @Description Admin endpoint to verify or reject a vendor's uploaded document
// @Tags Admin - Vendor Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Document ID (UUID)"
// @Param request body dto.AdminDocumentVerificationRequest true "Verification decision"
// @Success 200 {object} map[string]interface{} "Document verified successfully"
// @Failure 400 {object} apperrors.ErrorResponse "Invalid input or validation error"
// @Failure 401 {object} apperrors.ErrorResponse "Unauthorized"
// @Failure 403 {object} apperrors.ErrorResponse "Forbidden - Admin access required"
// @Failure 404 {object} apperrors.ErrorResponse "Document not found"
// @Failure 500 {object} apperrors.ErrorResponse "Internal server error"
// @Router /admin/vendors/documents/{id}/verify [put]
func (h *VendorHandler) VerifyDocument(c *gin.Context) {
	adminID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": apperrors.NewErrorResponse(apperrors.ErrUnauthorized),
		})
		return
	}

	documentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": apperrors.NewErrorResponse(apperrors.ErrInvalidInput),
		})
		return
	}

	var req dto.AdminDocumentVerificationRequest
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

	if err := h.vendorUseCase.VerifyDocument(c.Request.Context(), documentID, adminID, &req); err != nil {
		statusCode := apperrors.GetStatusCode(err)
		c.JSON(statusCode, gin.H{
			"error": apperrors.NewErrorResponse(err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Document verified successfully",
	})
}
