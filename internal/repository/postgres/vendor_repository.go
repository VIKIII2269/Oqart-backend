package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/oqart/backend/internal/domain"
	apperrors "github.com/oqart/backend/pkg/errors"
)

type VendorRepository struct {
	db *gorm.DB
}

func NewVendorRepository(db *gorm.DB) *VendorRepository {
	return &VendorRepository{db: db}
}

// Create creates a new vendor
func (r *VendorRepository) Create(ctx context.Context, vendor *domain.Vendor) error {
	if err := r.db.WithContext(ctx).Create(vendor).Error; err != nil {
		return apperrors.ErrDatabase
	}
	return nil
}

// GetByID retrieves a vendor by ID
func (r *VendorRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Vendor, error) {
	var vendor domain.Vendor
	if err := r.db.WithContext(ctx).
		Preload("Addresses").
		Preload("Documents").
		Preload("BankDetails").
		Preload("ContactPersons").
		Where("id = ?", id).
		First(&vendor).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrVendorNotFound
		}
		return nil, apperrors.ErrDatabase
	}
	return &vendor, nil
}

// GetByUserID retrieves a vendor by user ID
func (r *VendorRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Vendor, error) {
	var vendor domain.Vendor
	if err := r.db.WithContext(ctx).
		Preload("Addresses").
		Preload("Documents").
		Preload("BankDetails").
		Preload("ContactPersons").
		Where("user_id = ?", userID).
		First(&vendor).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrVendorNotFound
		}
		return nil, apperrors.ErrDatabase
	}
	return &vendor, nil
}

// Update updates an existing vendor
func (r *VendorRepository) Update(ctx context.Context, vendor *domain.Vendor) error {
	if err := r.db.WithContext(ctx).Save(vendor).Error; err != nil {
		return apperrors.ErrDatabase
	}
	return nil
}

// List retrieves vendors with pagination and filters
func (r *VendorRepository) List(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]*domain.Vendor, int64, error) {
	var vendors []*domain.Vendor
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Vendor{})

	// Apply filters
	if status, ok := filters["status"]; ok {
		query = query.Where("status = ?", status)
	}
	if isVerified, ok := filters["is_verified"]; ok {
		query = query.Where("is_verified = ?", isVerified)
	}
	if isFeatured, ok := filters["is_featured"]; ok {
		query = query.Where("is_featured = ?", isFeatured)
	}
	if search, ok := filters["search"]; ok {
		searchStr := "%" + search.(string) + "%"
		query = query.Where("business_name ILIKE ? OR description ILIKE ?", searchStr, searchStr)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, apperrors.ErrDatabase
	}

	// Get paginated results
	if err := query.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&vendors).Error; err != nil {
		return nil, 0, apperrors.ErrDatabase
	}

	return vendors, total, nil
}

// Exists checks if a vendor exists
func (r *VendorRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&domain.Vendor{}).
		Where("id = ?", id).
		Count(&count).Error; err != nil {
		return false, apperrors.ErrDatabase
	}
	return count > 0, nil
}

// ExistsByUserID checks if a vendor exists for a user
func (r *VendorRepository) ExistsByUserID(ctx context.Context, userID uuid.UUID) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&domain.Vendor{}).
		Where("user_id = ?", userID).
		Count(&count).Error; err != nil {
		return false, apperrors.ErrDatabase
	}
	return count > 0, nil
}

// Address operations

// CreateAddress creates a vendor address
func (r *VendorRepository) CreateAddress(ctx context.Context, address *domain.VendorAddress) error {
	if err := r.db.WithContext(ctx).Create(address).Error; err != nil {
		return apperrors.ErrDatabase
	}
	return nil
}

// GetAddress retrieves a vendor address by ID
func (r *VendorRepository) GetAddress(ctx context.Context, id uuid.UUID) (*domain.VendorAddress, error) {
	var address domain.VendorAddress
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&address).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, apperrors.ErrDatabase
	}
	return &address, nil
}

// UpdateAddress updates a vendor address
func (r *VendorRepository) UpdateAddress(ctx context.Context, address *domain.VendorAddress) error {
	if err := r.db.WithContext(ctx).Save(address).Error; err != nil {
		return apperrors.ErrDatabase
	}
	return nil
}

// SetPrimaryAddress sets an address as primary and unsets all others for the vendor
func (r *VendorRepository) SetPrimaryAddress(ctx context.Context, vendorID, addressID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Unset all primary addresses for the vendor
		if err := tx.Model(&domain.VendorAddress{}).
			Where("vendor_id = ? AND is_primary = ?", vendorID, true).
			Update("is_primary", false).Error; err != nil {
			return apperrors.ErrDatabase
		}

		// Set the new primary address
		if err := tx.Model(&domain.VendorAddress{}).
			Where("id = ? AND vendor_id = ?", addressID, vendorID).
			Update("is_primary", true).Error; err != nil {
			return apperrors.ErrDatabase
		}

		return nil
	})
}

// Bank details operations

// CreateBankDetails creates vendor bank details
func (r *VendorRepository) CreateBankDetails(ctx context.Context, bankDetails *domain.VendorBankDetails) error {
	if err := r.db.WithContext(ctx).Create(bankDetails).Error; err != nil {
		return apperrors.ErrDatabase
	}
	return nil
}

// GetBankDetails retrieves vendor bank details
func (r *VendorRepository) GetBankDetails(ctx context.Context, vendorID uuid.UUID) (*domain.VendorBankDetails, error) {
	var bankDetails domain.VendorBankDetails
	if err := r.db.WithContext(ctx).Where("vendor_id = ?", vendorID).First(&bankDetails).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, apperrors.ErrDatabase
	}
	return &bankDetails, nil
}

// UpdateBankDetails updates vendor bank details
func (r *VendorRepository) UpdateBankDetails(ctx context.Context, bankDetails *domain.VendorBankDetails) error {
	if err := r.db.WithContext(ctx).Save(bankDetails).Error; err != nil {
		return apperrors.ErrDatabase
	}
	return nil
}

// Document operations

// CreateDocument creates a vendor document
func (r *VendorRepository) CreateDocument(ctx context.Context, document *domain.VendorDocument) error {
	if err := r.db.WithContext(ctx).Create(document).Error; err != nil {
		return apperrors.ErrDatabase
	}
	return nil
}

// GetDocument retrieves a vendor document by ID
func (r *VendorRepository) GetDocument(ctx context.Context, id uuid.UUID) (*domain.VendorDocument, error) {
	var document domain.VendorDocument
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&document).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, apperrors.ErrDatabase
	}
	return &document, nil
}

// GetDocumentsByVendorID retrieves all documents for a vendor
func (r *VendorRepository) GetDocumentsByVendorID(ctx context.Context, vendorID uuid.UUID) ([]*domain.VendorDocument, error) {
	var documents []*domain.VendorDocument
	if err := r.db.WithContext(ctx).
		Where("vendor_id = ?", vendorID).
		Order("created_at DESC").
		Find(&documents).Error; err != nil {
		return nil, apperrors.ErrDatabase
	}
	return documents, nil
}

// UpdateDocument updates a vendor document
func (r *VendorRepository) UpdateDocument(ctx context.Context, document *domain.VendorDocument) error {
	if err := r.db.WithContext(ctx).Save(document).Error; err != nil {
		return apperrors.ErrDatabase
	}
	return nil
}

// DeleteDocument deletes a vendor document
func (r *VendorRepository) DeleteDocument(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&domain.VendorDocument{}, "id = ?", id).Error; err != nil {
		return apperrors.ErrDatabase
	}
	return nil
}

// Contact person operations

// CreateContactPerson creates a vendor contact person
func (r *VendorRepository) CreateContactPerson(ctx context.Context, contact *domain.VendorContactPerson) error {
	if err := r.db.WithContext(ctx).Create(contact).Error; err != nil {
		return apperrors.ErrDatabase
	}
	return nil
}

// GetContactPersonsByVendorID retrieves all contact persons for a vendor
func (r *VendorRepository) GetContactPersonsByVendorID(ctx context.Context, vendorID uuid.UUID) ([]*domain.VendorContactPerson, error) {
	var contacts []*domain.VendorContactPerson
	if err := r.db.WithContext(ctx).
		Where("vendor_id = ?", vendorID).
		Order("is_primary DESC, created_at DESC").
		Find(&contacts).Error; err != nil {
		return nil, apperrors.ErrDatabase
	}
	return contacts, nil
}

// SetPrimaryContactPerson sets a contact person as primary
func (r *VendorRepository) SetPrimaryContactPerson(ctx context.Context, vendorID, contactID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Unset all primary contacts for the vendor
		if err := tx.Model(&domain.VendorContactPerson{}).
			Where("vendor_id = ? AND is_primary = ?", vendorID, true).
			Update("is_primary", false).Error; err != nil {
			return apperrors.ErrDatabase
		}

		// Set the new primary contact
		if err := tx.Model(&domain.VendorContactPerson{}).
			Where("id = ? AND vendor_id = ?", contactID, vendorID).
			Update("is_primary", true).Error; err != nil {
			return apperrors.ErrDatabase
		}

		return nil
	})
}

// BelongsToVendor checks if a resource belongs to a vendor
func (r *VendorRepository) BelongsToVendor(ctx context.Context, vendorID uuid.UUID, resourceType string, resourceID uuid.UUID) (bool, error) {
	var count int64
	var err error

	switch resourceType {
	case "address":
		err = r.db.WithContext(ctx).Model(&domain.VendorAddress{}).
			Where("id = ? AND vendor_id = ?", resourceID, vendorID).
			Count(&count).Error
	case "document":
		err = r.db.WithContext(ctx).Model(&domain.VendorDocument{}).
			Where("id = ? AND vendor_id = ?", resourceID, vendorID).
			Count(&count).Error
	case "contact":
		err = r.db.WithContext(ctx).Model(&domain.VendorContactPerson{}).
			Where("id = ? AND vendor_id = ?", resourceID, vendorID).
			Count(&count).Error
	default:
		return false, apperrors.ErrInvalidInput
	}

	if err != nil {
		return false, apperrors.ErrDatabase
	}

	return count > 0, nil
}
