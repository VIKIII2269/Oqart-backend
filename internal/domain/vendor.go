package domain

import (
	"time"

	"github.com/google/uuid"
)

// Vendor status types
type VendorStatus string

const (
	VendorStatusPending      VendorStatus = "pending"
	VendorStatusUnderReview  VendorStatus = "under_review"
	VendorStatusApproved     VendorStatus = "approved"
	VendorStatusRejected     VendorStatus = "rejected"
	VendorStatusSuspended    VendorStatus = "suspended"
	VendorStatusInactive     VendorStatus = "inactive"
)

// Vendor business types
type VendorBusinessType string

const (
	BusinessTypeIndividual      VendorBusinessType = "individual"
	BusinessTypeProprietorship  VendorBusinessType = "proprietorship"
	BusinessTypePartnership     VendorBusinessType = "partnership"
	BusinessTypePrivateLimited  VendorBusinessType = "private_limited"
	BusinessTypePublicLimited   VendorBusinessType = "public_limited"
	BusinessTypeLLP             VendorBusinessType = "llp"
)

// Document types
type DocumentType string

const (
	DocumentTypeGSTCertificate    DocumentType = "gst_certificate"
	DocumentTypePANCard          DocumentType = "pan_card"
	DocumentTypeAddressProof     DocumentType = "address_proof"
	DocumentTypeBankStatement    DocumentType = "bank_statement"
	DocumentTypeCancelledCheque  DocumentType = "cancelled_cheque"
	DocumentTypeFSSAILicense     DocumentType = "fssai_license"
	DocumentTypeOrganicCertificate DocumentType = "organic_certificate"
)

// Document status
type DocumentStatus string

const (
	DocumentStatusPending      DocumentStatus = "pending"
	DocumentStatusUnderReview  DocumentStatus = "under_review"
	DocumentStatusApproved     DocumentStatus = "approved"
	DocumentStatusRejected     DocumentStatus = "rejected"
)

// Vendor represents a vendor/seller in the marketplace
type Vendor struct {
	ID              uuid.UUID          `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID          uuid.UUID          `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`
	BusinessName    string             `gorm:"type:varchar(255);not null" json:"business_name"`
	BusinessType    VendorBusinessType `gorm:"type:varchar(50);not null" json:"business_type"`
	GSTIN           *string            `gorm:"type:varchar(15);uniqueIndex" json:"gstin,omitempty"`
	PAN             *string            `gorm:"type:varchar(10)" json:"pan,omitempty"`
	FSSAINumber     *string            `gorm:"type:varchar(14)" json:"fssai_number,omitempty"`
	Status          VendorStatus       `gorm:"type:varchar(50);not null;default:'pending'" json:"status"`
	OnboardingStep  int                `gorm:"default:1" json:"onboarding_step"`
	RejectionReason *string            `gorm:"type:text" json:"rejection_reason,omitempty"`
	ApprovedAt      *time.Time         `json:"approved_at,omitempty"`
	ApprovedBy      *uuid.UUID         `gorm:"type:uuid" json:"approved_by,omitempty"`
	Rating          float64            `gorm:"type:decimal(3,2);default:0.00" json:"rating"`
	TotalReviews    int                `gorm:"default:0" json:"total_reviews"`
	TotalOrders     int                `gorm:"default:0" json:"total_orders"`
	TotalSales      float64            `gorm:"type:decimal(15,2);default:0.00" json:"total_sales"`
	CommissionRate  float64            `gorm:"type:decimal(5,2);default:10.00" json:"commission_rate"`
	IsVerified      bool               `gorm:"default:false" json:"is_verified"`
	IsFeatured      bool               `gorm:"default:false" json:"is_featured"`
	Description     *string            `gorm:"type:text" json:"description,omitempty"`
	LogoURL         *string            `gorm:"type:varchar(500)" json:"logo_url,omitempty"`
	BannerURL       *string            `gorm:"type:varchar(500)" json:"banner_url,omitempty"`
	CreatedAt       time.Time          `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time          `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`

	// Relationships
	User             *User                  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Addresses        []VendorAddress        `gorm:"foreignKey:VendorID" json:"addresses,omitempty"`
	Documents        []VendorDocument       `gorm:"foreignKey:VendorID" json:"documents,omitempty"`
	BankDetails      *VendorBankDetails     `gorm:"foreignKey:VendorID" json:"bank_details,omitempty"`
	ContactPersons   []VendorContactPerson  `gorm:"foreignKey:VendorID" json:"contact_persons,omitempty"`
	BusinessHours    []VendorBusinessHours  `gorm:"foreignKey:VendorID" json:"business_hours,omitempty"`
}

// VendorAddress represents a vendor's business address
type VendorAddress struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	VendorID      uuid.UUID `gorm:"type:uuid;not null;index" json:"vendor_id"`
	Type          string    `gorm:"type:varchar(50);not null;default:'business'" json:"type"`
	StreetAddress string    `gorm:"type:varchar(255);not null" json:"street_address"`
	Landmark      *string   `gorm:"type:varchar(255)" json:"landmark,omitempty"`
	City          string    `gorm:"type:varchar(100);not null" json:"city"`
	State         string    `gorm:"type:varchar(100);not null" json:"state"`
	Pincode       string    `gorm:"type:varchar(10);not null" json:"pincode"`
	Country       string    `gorm:"type:varchar(100);default:'India'" json:"country"`
	Latitude      *float64  `gorm:"type:decimal(10,8)" json:"latitude,omitempty"`
	Longitude     *float64  `gorm:"type:decimal(11,8)" json:"longitude,omitempty"`
	IsPrimary     bool      `gorm:"default:false" json:"is_primary"`
	CreatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// VendorDocument represents uploaded vendor documents
type VendorDocument struct {
	ID             uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	VendorID       uuid.UUID      `gorm:"type:uuid;not null;index" json:"vendor_id"`
	Type           DocumentType   `gorm:"type:varchar(50);not null" json:"type"`
	FilePath       string         `gorm:"type:varchar(500);not null" json:"file_path"`
	FileName       string         `gorm:"type:varchar(255);not null" json:"file_name"`
	FileSize       int            `gorm:"not null" json:"file_size"`
	MimeType       string         `gorm:"type:varchar(100);not null" json:"mime_type"`
	Status         DocumentStatus `gorm:"type:varchar(50);not null;default:'pending'" json:"status"`
	ExtractedData  *string        `gorm:"type:jsonb" json:"extracted_data,omitempty"`
	RejectionReason *string       `gorm:"type:text" json:"rejection_reason,omitempty"`
	VerifiedBy     *uuid.UUID     `gorm:"type:uuid" json:"verified_by,omitempty"`
	VerifiedAt     *time.Time     `json:"verified_at,omitempty"`
	CreatedAt      time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// VendorBankDetails represents vendor's bank account information
type VendorBankDetails struct {
	ID                uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	VendorID          uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex" json:"vendor_id"`
	AccountHolderName string     `gorm:"type:varchar(255);not null" json:"account_holder_name"`
	AccountNumber     string     `gorm:"type:varchar(50);not null" json:"account_number"`
	IFSCCode          string     `gorm:"type:varchar(11);not null" json:"ifsc_code"`
	BankName          string     `gorm:"type:varchar(255);not null" json:"bank_name"`
	BranchName        *string    `gorm:"type:varchar(255)" json:"branch_name,omitempty"`
	AccountType       string     `gorm:"type:varchar(50);default:'savings'" json:"account_type"`
	IsVerified        bool       `gorm:"default:false" json:"is_verified"`
	VerifiedAt        *time.Time `json:"verified_at,omitempty"`
	CreatedAt         time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt         time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// VendorContactPerson represents vendor contact persons
type VendorContactPerson struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	VendorID    uuid.UUID `gorm:"type:uuid;not null;index" json:"vendor_id"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Designation *string   `gorm:"type:varchar(100)" json:"designation,omitempty"`
	Email       *string   `gorm:"type:varchar(255)" json:"email,omitempty"`
	Phone       string    `gorm:"type:varchar(20);not null" json:"phone"`
	IsPrimary   bool      `gorm:"default:false" json:"is_primary"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// VendorBusinessHours represents vendor's operating hours
type VendorBusinessHours struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	VendorID    uuid.UUID `gorm:"type:uuid;not null;index" json:"vendor_id"`
	DayOfWeek   int       `gorm:"not null;check:day_of_week >= 0 AND day_of_week <= 6" json:"day_of_week"`
	IsOpen      bool      `gorm:"default:true" json:"is_open"`
	OpeningTime *string   `gorm:"type:time" json:"opening_time,omitempty"`
	ClosingTime *string   `gorm:"type:time" json:"closing_time,omitempty"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName overrides
func (Vendor) TableName() string {
	return "vendors"
}

func (VendorAddress) TableName() string {
	return "vendor_addresses"
}

func (VendorDocument) TableName() string {
	return "vendor_documents"
}

func (VendorBankDetails) TableName() string {
	return "vendor_bank_details"
}

func (VendorContactPerson) TableName() string {
	return "vendor_contact_persons"
}

func (VendorBusinessHours) TableName() string {
	return "vendor_business_hours"
}
