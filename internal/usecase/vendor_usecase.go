package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/oqart/backend/internal/delivery/http/dto"
	"github.com/oqart/backend/internal/domain"
	"github.com/oqart/backend/internal/repository/postgres"
	"github.com/oqart/backend/pkg/cache"
	apperrors "github.com/oqart/backend/pkg/errors"
	"github.com/oqart/backend/pkg/logger"
)

type VendorUseCase interface {
	// Onboarding
	OnboardStep1(ctx context.Context, userID uuid.UUID, req *dto.VendorOnboardingStep1Request) (*dto.VendorResponse, error)
	OnboardStep2(ctx context.Context, userID uuid.UUID, req *dto.VendorOnboardingStep2Request) (*dto.VendorAddressResponse, error)
	OnboardStep3(ctx context.Context, userID uuid.UUID, req *dto.VendorOnboardingStep3Request) (*dto.VendorBankDetailsResponse, error)
	OnboardStep4(ctx context.Context, userID uuid.UUID, req *dto.VendorOnboardingStep4Request) (*dto.VendorContactPersonResponse, error)

	// Vendor profile
	GetMyVendor(ctx context.Context, userID uuid.UUID) (*dto.VendorDetailedResponse, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, req *dto.UpdateVendorProfileRequest) (*dto.VendorResponse, error)

	// Public vendor info
	GetVendorByID(ctx context.Context, vendorID uuid.UUID) (*dto.VendorResponse, error)
	ListVendors(ctx context.Context, filters map[string]interface{}, page, pageSize int) (*dto.VendorListResponse, error)

	// Documents
	GetDocuments(ctx context.Context, userID uuid.UUID) ([]dto.VendorDocumentResponse, error)

	// Admin operations
	ApproveVendor(ctx context.Context, vendorID uuid.UUID, adminID uuid.UUID, req *dto.AdminVendorApprovalRequest) error
	VerifyDocument(ctx context.Context, documentID uuid.UUID, adminID uuid.UUID, req *dto.AdminDocumentVerificationRequest) error
}

type vendorUseCase struct {
	vendorRepo *postgres.VendorRepository
	userRepo   postgres.UserRepository
	cache      *cache.Cache
	logger     *logger.Logger
}

func NewVendorUseCase(
	vendorRepo *postgres.VendorRepository,
	userRepo postgres.UserRepository,
	cache *cache.Cache,
	logger *logger.Logger,
) VendorUseCase {
	return &vendorUseCase{
		vendorRepo: vendorRepo,
		userRepo:   userRepo,
		cache:      cache,
		logger:     logger,
	}
}

// OnboardStep1 handles the first step of vendor onboarding - business details
func (uc *vendorUseCase) OnboardStep1(ctx context.Context, userID uuid.UUID, req *dto.VendorOnboardingStep1Request) (*dto.VendorResponse, error) {
	// Check if user already has a vendor account
	exists, err := uc.vendorRepo.ExistsByUserID(ctx, userID)
	if err != nil {
		uc.logger.Error("Failed to check vendor existence", "error", err, "user_id", userID)
		return nil, err
	}
	if exists {
		return nil, apperrors.New("VENDOR_ALREADY_EXISTS", "Vendor account already exists for this user", 409)
	}

	// Validate GSTIN if provided
	if req.GSTIN != nil && *req.GSTIN != "" {
		*req.GSTIN = strings.ToUpper(strings.TrimSpace(*req.GSTIN))
	}

	// Validate PAN if provided
	if req.PAN != nil && *req.PAN != "" {
		*req.PAN = strings.ToUpper(strings.TrimSpace(*req.PAN))
	}

	vendor := &domain.Vendor{
		ID:             uuid.New(),
		UserID:         userID,
		BusinessName:   req.BusinessName,
		BusinessType:   domain.VendorBusinessType(req.BusinessType),
		GSTIN:          req.GSTIN,
		PAN:            req.PAN,
		FSSAINumber:    req.FSSAINumber,
		Status:         domain.VendorStatusPending,
		OnboardingStep: 1,
		Description:    req.Description,
	}

	if err := uc.vendorRepo.Create(ctx, vendor); err != nil {
		uc.logger.Error("Failed to create vendor", "error", err, "user_id", userID)
		return nil, err
	}

	uc.logger.Info("Vendor onboarding step 1 completed", "vendor_id", vendor.ID, "user_id", userID)
	return uc.mapVendorToResponse(vendor), nil
}

// OnboardStep2 handles the second step - business address
func (uc *vendorUseCase) OnboardStep2(ctx context.Context, userID uuid.UUID, req *dto.VendorOnboardingStep2Request) (*dto.VendorAddressResponse, error) {
	vendor, err := uc.vendorRepo.GetByUserID(ctx, userID)
	if err != nil {
		uc.logger.Error("Failed to get vendor", "error", err, "user_id", userID)
		return nil, err
	}

	if vendor.OnboardingStep < 1 {
		return nil, apperrors.New("INVALID_ONBOARDING_STEP", "Complete step 1 first", 400)
	}

	address := &domain.VendorAddress{
		ID:            uuid.New(),
		VendorID:      vendor.ID,
		Type:          "business",
		StreetAddress: req.StreetAddress,
		Landmark:      req.Landmark,
		City:          req.City,
		State:         req.State,
		Pincode:       req.Pincode,
		Country:       req.Country,
		Latitude:      req.Latitude,
		Longitude:     req.Longitude,
		IsPrimary:     true,
	}

	if err := uc.vendorRepo.CreateAddress(ctx, address); err != nil {
		uc.logger.Error("Failed to create vendor address", "error", err, "vendor_id", vendor.ID)
		return nil, err
	}

	// Update onboarding step
	vendor.OnboardingStep = 2
	if err := uc.vendorRepo.Update(ctx, vendor); err != nil {
		uc.logger.Error("Failed to update vendor onboarding step", "error", err, "vendor_id", vendor.ID)
		return nil, err
	}

	uc.logger.Info("Vendor onboarding step 2 completed", "vendor_id", vendor.ID)
	return uc.mapAddressToResponse(address), nil
}

// OnboardStep3 handles the third step - bank details
func (uc *vendorUseCase) OnboardStep3(ctx context.Context, userID uuid.UUID, req *dto.VendorOnboardingStep3Request) (*dto.VendorBankDetailsResponse, error) {
	vendor, err := uc.vendorRepo.GetByUserID(ctx, userID)
	if err != nil {
		uc.logger.Error("Failed to get vendor", "error", err, "user_id", userID)
		return nil, err
	}

	if vendor.OnboardingStep < 2 {
		return nil, apperrors.New("INVALID_ONBOARDING_STEP", "Complete step 2 first", 400)
	}

	bankDetails := &domain.VendorBankDetails{
		ID:                uuid.New(),
		VendorID:          vendor.ID,
		AccountHolderName: req.AccountHolderName,
		AccountNumber:     req.AccountNumber,
		IFSCCode:          strings.ToUpper(req.IFSCCode),
		BankName:          req.BankName,
		BranchName:        req.BranchName,
		AccountType:       req.AccountType,
		IsVerified:        false,
	}

	if err := uc.vendorRepo.CreateBankDetails(ctx, bankDetails); err != nil {
		uc.logger.Error("Failed to create bank details", "error", err, "vendor_id", vendor.ID)
		return nil, err
	}

	// Update onboarding step
	vendor.OnboardingStep = 3
	if err := uc.vendorRepo.Update(ctx, vendor); err != nil {
		uc.logger.Error("Failed to update vendor onboarding step", "error", err, "vendor_id", vendor.ID)
		return nil, err
	}

	uc.logger.Info("Vendor onboarding step 3 completed", "vendor_id", vendor.ID)
	return uc.mapBankDetailsToResponse(bankDetails), nil
}

// OnboardStep4 handles the fourth step - contact person
func (uc *vendorUseCase) OnboardStep4(ctx context.Context, userID uuid.UUID, req *dto.VendorOnboardingStep4Request) (*dto.VendorContactPersonResponse, error) {
	vendor, err := uc.vendorRepo.GetByUserID(ctx, userID)
	if err != nil {
		uc.logger.Error("Failed to get vendor", "error", err, "user_id", userID)
		return nil, err
	}

	if vendor.OnboardingStep < 3 {
		return nil, apperrors.New("INVALID_ONBOARDING_STEP", "Complete step 3 first", 400)
	}

	contact := &domain.VendorContactPerson{
		ID:          uuid.New(),
		VendorID:    vendor.ID,
		Name:        req.Name,
		Designation: req.Designation,
		Email:       req.Email,
		Phone:       req.Phone,
		IsPrimary:   true,
	}

	if err := uc.vendorRepo.CreateContactPerson(ctx, contact); err != nil {
		uc.logger.Error("Failed to create contact person", "error", err, "vendor_id", vendor.ID)
		return nil, err
	}

	// Update onboarding step and status
	vendor.OnboardingStep = 4
	vendor.Status = domain.VendorStatusUnderReview
	if err := uc.vendorRepo.Update(ctx, vendor); err != nil {
		uc.logger.Error("Failed to update vendor status", "error", err, "vendor_id", vendor.ID)
		return nil, err
	}

	uc.logger.Info("Vendor onboarding completed", "vendor_id", vendor.ID, "status", vendor.Status)
	return uc.mapContactPersonToResponse(contact), nil
}

// GetMyVendor retrieves the current user's vendor profile
func (uc *vendorUseCase) GetMyVendor(ctx context.Context, userID uuid.UUID) (*dto.VendorDetailedResponse, error) {
	vendor, err := uc.vendorRepo.GetByUserID(ctx, userID)
	if err != nil {
		uc.logger.Error("Failed to get vendor", "error", err, "user_id", userID)
		return nil, err
	}

	response := &dto.VendorDetailedResponse{
		Vendor: *uc.mapVendorToResponse(vendor),
	}

	// Add address if exists
	if len(vendor.Addresses) > 0 {
		for _, addr := range vendor.Addresses {
			if addr.IsPrimary {
				response.Address = uc.mapAddressToResponse(&addr)
				break
			}
		}
		if response.Address == nil {
			response.Address = uc.mapAddressToResponse(&vendor.Addresses[0])
		}
	}

	// Add bank details if exists
	if vendor.BankDetails != nil {
		response.BankDetails = uc.mapBankDetailsToResponse(vendor.BankDetails)
	}

	// Add primary contact if exists
	if len(vendor.ContactPersons) > 0 {
		for _, contact := range vendor.ContactPersons {
			if contact.IsPrimary {
				response.ContactPerson = uc.mapContactPersonToResponse(&contact)
				break
			}
		}
		if response.ContactPerson == nil {
			response.ContactPerson = uc.mapContactPersonToResponse(&vendor.ContactPersons[0])
		}
	}

	// Add documents
	if len(vendor.Documents) > 0 {
		documents := make([]dto.VendorDocumentResponse, len(vendor.Documents))
		for i, doc := range vendor.Documents {
			documents[i] = *uc.mapDocumentToResponse(&doc)
		}
		response.Documents = documents
	}

	return response, nil
}

// UpdateProfile updates vendor profile
func (uc *vendorUseCase) UpdateProfile(ctx context.Context, userID uuid.UUID, req *dto.UpdateVendorProfileRequest) (*dto.VendorResponse, error) {
	vendor, err := uc.vendorRepo.GetByUserID(ctx, userID)
	if err != nil {
		uc.logger.Error("Failed to get vendor", "error", err, "user_id", userID)
		return nil, err
	}

	if req.BusinessName != nil {
		vendor.BusinessName = *req.BusinessName
	}
	if req.Description != nil {
		vendor.Description = req.Description
	}
	if req.LogoURL != nil {
		vendor.LogoURL = req.LogoURL
	}
	if req.BannerURL != nil {
		vendor.BannerURL = req.BannerURL
	}

	if err := uc.vendorRepo.Update(ctx, vendor); err != nil {
		uc.logger.Error("Failed to update vendor profile", "error", err, "vendor_id", vendor.ID)
		return nil, err
	}

	uc.logger.Info("Vendor profile updated", "vendor_id", vendor.ID)
	return uc.mapVendorToResponse(vendor), nil
}

// GetVendorByID retrieves a vendor by ID (public)
func (uc *vendorUseCase) GetVendorByID(ctx context.Context, vendorID uuid.UUID) (*dto.VendorResponse, error) {
	vendor, err := uc.vendorRepo.GetByID(ctx, vendorID)
	if err != nil {
		return nil, err
	}

	// Only return approved vendors for public view
	if vendor.Status != domain.VendorStatusApproved {
		return nil, apperrors.ErrVendorNotFound
	}

	return uc.mapVendorToResponse(vendor), nil
}

// ListVendors retrieves a list of vendors with pagination
func (uc *vendorUseCase) ListVendors(ctx context.Context, filters map[string]interface{}, page, pageSize int) (*dto.VendorListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	vendors, total, err := uc.vendorRepo.List(ctx, filters, pageSize, offset)
	if err != nil {
		uc.logger.Error("Failed to list vendors", "error", err)
		return nil, err
	}

	responses := make([]dto.VendorResponse, len(vendors))
	for i, vendor := range vendors {
		responses[i] = *uc.mapVendorToResponse(vendor)
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &dto.VendorListResponse{
		Vendors:    responses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// GetDocuments retrieves all documents for a vendor
func (uc *vendorUseCase) GetDocuments(ctx context.Context, userID uuid.UUID) ([]dto.VendorDocumentResponse, error) {
	vendor, err := uc.vendorRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	documents, err := uc.vendorRepo.GetDocumentsByVendorID(ctx, vendor.ID)
	if err != nil {
		uc.logger.Error("Failed to get documents", "error", err, "vendor_id", vendor.ID)
		return nil, err
	}

	responses := make([]dto.VendorDocumentResponse, len(documents))
	for i, doc := range documents {
		responses[i] = *uc.mapDocumentToResponse(doc)
	}

	return responses, nil
}

// ApproveVendor approves or rejects a vendor (admin only)
func (uc *vendorUseCase) ApproveVendor(ctx context.Context, vendorID uuid.UUID, adminID uuid.UUID, req *dto.AdminVendorApprovalRequest) error {
	vendor, err := uc.vendorRepo.GetByID(ctx, vendorID)
	if err != nil {
		return err
	}

	switch req.Status {
	case "approved":
		vendor.Status = domain.VendorStatusApproved
		now := new(time.Time)
		*now = time.Now()
		vendor.ApprovedAt = now
		vendor.ApprovedBy = &adminID
		vendor.RejectionReason = nil
	case "rejected":
		vendor.Status = domain.VendorStatusRejected
		vendor.RejectionReason = req.RejectionReason
	case "suspended":
		vendor.Status = domain.VendorStatusSuspended
		vendor.RejectionReason = req.RejectionReason
	default:
		return apperrors.ErrInvalidInput
	}

	if err := uc.vendorRepo.Update(ctx, vendor); err != nil {
		uc.logger.Error("Failed to update vendor status", "error", err, "vendor_id", vendorID)
		return err
	}

	uc.logger.Info("Vendor status updated by admin", "vendor_id", vendorID, "status", vendor.Status, "admin_id", adminID)
	return nil
}

// VerifyDocument verifies or rejects a document (admin only)
func (uc *vendorUseCase) VerifyDocument(ctx context.Context, documentID uuid.UUID, adminID uuid.UUID, req *dto.AdminDocumentVerificationRequest) error {
	document, err := uc.vendorRepo.GetDocument(ctx, documentID)
	if err != nil {
		return err
	}

	switch req.Status {
	case "approved":
		document.Status = domain.DocumentStatusApproved
		now := new(time.Time)
		*now = time.Now()
		document.VerifiedAt = now
		document.VerifiedBy = &adminID
		document.RejectionReason = nil
	case "rejected":
		document.Status = domain.DocumentStatusRejected
		document.RejectionReason = req.RejectionReason
	default:
		return apperrors.ErrInvalidInput
	}

	if err := uc.vendorRepo.UpdateDocument(ctx, document); err != nil {
		uc.logger.Error("Failed to update document status", "error", err, "document_id", documentID)
		return err
	}

	uc.logger.Info("Document verified by admin", "document_id", documentID, "status", document.Status, "admin_id", adminID)
	return nil
}

// Helper methods

func (uc *vendorUseCase) mapVendorToResponse(vendor *domain.Vendor) *dto.VendorResponse {
	return &dto.VendorResponse{
		ID:              vendor.ID.String(),
		UserID:          vendor.UserID.String(),
		BusinessName:    vendor.BusinessName,
		BusinessType:    string(vendor.BusinessType),
		GSTIN:           vendor.GSTIN,
		PAN:             vendor.PAN,
		FSSAINumber:     vendor.FSSAINumber,
		Status:          string(vendor.Status),
		OnboardingStep:  vendor.OnboardingStep,
		RejectionReason: vendor.RejectionReason,
		ApprovedAt:      vendor.ApprovedAt,
		Rating:          vendor.Rating,
		TotalReviews:    vendor.TotalReviews,
		TotalOrders:     vendor.TotalOrders,
		TotalSales:      vendor.TotalSales,
		CommissionRate:  vendor.CommissionRate,
		IsVerified:      vendor.IsVerified,
		IsFeatured:      vendor.IsFeatured,
		Description:     vendor.Description,
		LogoURL:         vendor.LogoURL,
		BannerURL:       vendor.BannerURL,
		CreatedAt:       vendor.CreatedAt,
		UpdatedAt:       vendor.UpdatedAt,
	}
}

func (uc *vendorUseCase) mapAddressToResponse(addr *domain.VendorAddress) *dto.VendorAddressResponse {
	return &dto.VendorAddressResponse{
		ID:            addr.ID.String(),
		VendorID:      addr.VendorID.String(),
		Type:          addr.Type,
		StreetAddress: addr.StreetAddress,
		Landmark:      addr.Landmark,
		City:          addr.City,
		State:         addr.State,
		Pincode:       addr.Pincode,
		Country:       addr.Country,
		Latitude:      addr.Latitude,
		Longitude:     addr.Longitude,
		IsPrimary:     addr.IsPrimary,
		CreatedAt:     addr.CreatedAt,
		UpdatedAt:     addr.UpdatedAt,
	}
}

func (uc *vendorUseCase) mapBankDetailsToResponse(bank *domain.VendorBankDetails) *dto.VendorBankDetailsResponse {
	// Mask account number for security
	maskedAccount := "XXXX" + bank.AccountNumber[len(bank.AccountNumber)-4:]

	return &dto.VendorBankDetailsResponse{
		ID:                bank.ID.String(),
		VendorID:          bank.VendorID.String(),
		AccountHolderName: bank.AccountHolderName,
		AccountNumber:     maskedAccount,
		IFSCCode:          bank.IFSCCode,
		BankName:          bank.BankName,
		BranchName:        bank.BranchName,
		AccountType:       bank.AccountType,
		IsVerified:        bank.IsVerified,
		VerifiedAt:        bank.VerifiedAt,
		CreatedAt:         bank.CreatedAt,
		UpdatedAt:         bank.UpdatedAt,
	}
}

func (uc *vendorUseCase) mapContactPersonToResponse(contact *domain.VendorContactPerson) *dto.VendorContactPersonResponse {
	return &dto.VendorContactPersonResponse{
		ID:          contact.ID.String(),
		VendorID:    contact.VendorID.String(),
		Name:        contact.Name,
		Designation: contact.Designation,
		Email:       contact.Email,
		Phone:       contact.Phone,
		IsPrimary:   contact.IsPrimary,
		CreatedAt:   contact.CreatedAt,
		UpdatedAt:   contact.UpdatedAt,
	}
}

func (uc *vendorUseCase) mapDocumentToResponse(doc *domain.VendorDocument) *dto.VendorDocumentResponse {
	return &dto.VendorDocumentResponse{
		ID:              doc.ID.String(),
		VendorID:        doc.VendorID.String(),
		Type:            string(doc.Type),
		FileName:        doc.FileName,
		FileSize:        doc.FileSize,
		Status:          string(doc.Status),
		RejectionReason: doc.RejectionReason,
		VerifiedAt:      doc.VerifiedAt,
		CreatedAt:       doc.CreatedAt,
		UpdatedAt:       doc.UpdatedAt,
	}
}
