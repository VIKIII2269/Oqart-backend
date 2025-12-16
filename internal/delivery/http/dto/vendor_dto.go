package dto

import "time"

// VendorOnboardingStep1Request represents the first step of vendor onboarding
type VendorOnboardingStep1Request struct {
	BusinessName string `json:"business_name" validate:"required,min=3,max=255"`
	BusinessType string `json:"business_type" validate:"required,oneof=individual proprietorship partnership private_limited public_limited llp"`
	GSTIN        *string `json:"gstin,omitempty" validate:"omitempty,gstin"`
	PAN          *string `json:"pan,omitempty" validate:"omitempty,pan"`
	FSSAINumber  *string `json:"fssai_number,omitempty" validate:"omitempty,len=14"`
	Description  *string `json:"description,omitempty" validate:"omitempty,max=1000"`
}

// VendorOnboardingStep2Request represents address information
type VendorOnboardingStep2Request struct {
	StreetAddress string   `json:"street_address" validate:"required,min=5,max=255"`
	Landmark      *string  `json:"landmark,omitempty" validate:"omitempty,max=255"`
	City          string   `json:"city" validate:"required,min=2,max=100"`
	State         string   `json:"state" validate:"required,min=2,max=100"`
	Pincode       string   `json:"pincode" validate:"required,indian_pincode"`
	Country       string   `json:"country" validate:"required,min=2,max=100"`
	Latitude      *float64 `json:"latitude,omitempty" validate:"omitempty,latitude"`
	Longitude     *float64 `json:"longitude,omitempty" validate:"omitempty,longitude"`
}

// VendorOnboardingStep3Request represents bank details
type VendorOnboardingStep3Request struct {
	AccountHolderName string  `json:"account_holder_name" validate:"required,min=3,max=255"`
	AccountNumber     string  `json:"account_number" validate:"required,min=9,max=18"`
	IFSCCode          string  `json:"ifsc_code" validate:"required,ifsc"`
	BankName          string  `json:"bank_name" validate:"required,min=2,max=255"`
	BranchName        *string `json:"branch_name,omitempty" validate:"omitempty,max=255"`
	AccountType       string  `json:"account_type" validate:"required,oneof=savings current"`
}

// VendorOnboardingStep4Request represents contact person information
type VendorOnboardingStep4Request struct {
	Name        string  `json:"name" validate:"required,min=3,max=255"`
	Designation *string `json:"designation,omitempty" validate:"omitempty,max=100"`
	Email       *string `json:"email,omitempty" validate:"omitempty,email"`
	Phone       string  `json:"phone" validate:"required,indian_phone"`
}

// VendorDocumentUploadRequest represents document upload
type VendorDocumentUploadRequest struct {
	DocumentType string `json:"document_type" validate:"required,oneof=gst_certificate pan_card address_proof bank_statement cancelled_cheque fssai_license organic_certificate"`
}

// VendorResponse represents vendor information
type VendorResponse struct {
	ID              string     `json:"id"`
	UserID          string     `json:"user_id"`
	BusinessName    string     `json:"business_name"`
	BusinessType    string     `json:"business_type"`
	GSTIN           *string    `json:"gstin,omitempty"`
	PAN             *string    `json:"pan,omitempty"`
	FSSAINumber     *string    `json:"fssai_number,omitempty"`
	Status          string     `json:"status"`
	OnboardingStep  int        `json:"onboarding_step"`
	RejectionReason *string    `json:"rejection_reason,omitempty"`
	ApprovedAt      *time.Time `json:"approved_at,omitempty"`
	Rating          float64    `json:"rating"`
	TotalReviews    int        `json:"total_reviews"`
	TotalOrders     int        `json:"total_orders"`
	TotalSales      float64    `json:"total_sales"`
	CommissionRate  float64    `json:"commission_rate"`
	IsVerified      bool       `json:"is_verified"`
	IsFeatured      bool       `json:"is_featured"`
	Description     *string    `json:"description,omitempty"`
	LogoURL         *string    `json:"logo_url,omitempty"`
	BannerURL       *string    `json:"banner_url,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// VendorAddressResponse represents vendor address
type VendorAddressResponse struct {
	ID            string    `json:"id"`
	VendorID      string    `json:"vendor_id"`
	Type          string    `json:"type"`
	StreetAddress string    `json:"street_address"`
	Landmark      *string   `json:"landmark,omitempty"`
	City          string    `json:"city"`
	State         string    `json:"state"`
	Pincode       string    `json:"pincode"`
	Country       string    `json:"country"`
	Latitude      *float64  `json:"latitude,omitempty"`
	Longitude     *float64  `json:"longitude,omitempty"`
	IsPrimary     bool      `json:"is_primary"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// VendorDocumentResponse represents vendor document
type VendorDocumentResponse struct {
	ID              string     `json:"id"`
	VendorID        string     `json:"vendor_id"`
	Type            string     `json:"type"`
	FileName        string     `json:"file_name"`
	FileSize        int        `json:"file_size"`
	Status          string     `json:"status"`
	RejectionReason *string    `json:"rejection_reason,omitempty"`
	VerifiedAt      *time.Time `json:"verified_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// VendorBankDetailsResponse represents vendor bank details
type VendorBankDetailsResponse struct {
	ID                string     `json:"id"`
	VendorID          string     `json:"vendor_id"`
	AccountHolderName string     `json:"account_holder_name"`
	AccountNumber     string     `json:"account_number_masked"` // Masked for security
	IFSCCode          string     `json:"ifsc_code"`
	BankName          string     `json:"bank_name"`
	BranchName        *string    `json:"branch_name,omitempty"`
	AccountType       string     `json:"account_type"`
	IsVerified        bool       `json:"is_verified"`
	VerifiedAt        *time.Time `json:"verified_at,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

// VendorContactPersonResponse represents vendor contact person
type VendorContactPersonResponse struct {
	ID          string    `json:"id"`
	VendorID    string    `json:"vendor_id"`
	Name        string    `json:"name"`
	Designation *string   `json:"designation,omitempty"`
	Email       *string   `json:"email,omitempty"`
	Phone       string    `json:"phone"`
	IsPrimary   bool      `json:"is_primary"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// VendorDetailedResponse represents complete vendor information
type VendorDetailedResponse struct {
	Vendor         VendorResponse                `json:"vendor"`
	Address        *VendorAddressResponse        `json:"address,omitempty"`
	BankDetails    *VendorBankDetailsResponse    `json:"bank_details,omitempty"`
	ContactPerson  *VendorContactPersonResponse  `json:"contact_person,omitempty"`
	Documents      []VendorDocumentResponse      `json:"documents,omitempty"`
}

// UpdateVendorProfileRequest represents vendor profile update
type UpdateVendorProfileRequest struct {
	BusinessName *string  `json:"business_name,omitempty" validate:"omitempty,min=3,max=255"`
	Description  *string  `json:"description,omitempty" validate:"omitempty,max=1000"`
	LogoURL      *string  `json:"logo_url,omitempty" validate:"omitempty,url"`
	BannerURL    *string  `json:"banner_url,omitempty" validate:"omitempty,url"`
}

// VendorListResponse represents a list of vendors
type VendorListResponse struct {
	Vendors    []VendorResponse `json:"vendors"`
	Total      int64            `json:"total"`
	Page       int              `json:"page"`
	PageSize   int              `json:"page_size"`
	TotalPages int              `json:"total_pages"`
}

// AdminVendorApprovalRequest represents admin approval/rejection
type AdminVendorApprovalRequest struct {
	Status         string  `json:"status" validate:"required,oneof=approved rejected suspended"`
	RejectionReason *string `json:"rejection_reason,omitempty" validate:"required_if=Status rejected"`
}

// AdminDocumentVerificationRequest represents document verification
type AdminDocumentVerificationRequest struct {
	Status          string  `json:"status" validate:"required,oneof=approved rejected"`
	RejectionReason *string `json:"rejection_reason,omitempty" validate:"required_if=Status rejected"`
}
