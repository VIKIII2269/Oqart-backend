package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/oqart/backend/internal/domain"
	apperrors "github.com/oqart/backend/pkg/errors"
)

type AddressRepository struct {
	db *gorm.DB
}

func NewAddressRepository(db *gorm.DB) *AddressRepository {
	return &AddressRepository{db: db}
}

// Create creates a new address
func (r *AddressRepository) Create(ctx context.Context, address *domain.Address) error {
	if err := r.db.WithContext(ctx).Create(address).Error; err != nil {
		return apperrors.ErrDatabase
	}
	return nil
}

// GetByID retrieves an address by ID
func (r *AddressRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Address, error) {
	var address domain.Address
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&address).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrAddressNotFound
		}
		return nil, apperrors.ErrDatabase
	}
	return &address, nil
}

// GetByUserID retrieves all addresses for a user
func (r *AddressRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Address, error) {
	var addresses []*domain.Address
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("is_default DESC, created_at DESC").
		Find(&addresses).Error; err != nil {
		return nil, apperrors.ErrDatabase
	}
	return addresses, nil
}

// Update updates an existing address
func (r *AddressRepository) Update(ctx context.Context, address *domain.Address) error {
	if err := r.db.WithContext(ctx).Save(address).Error; err != nil {
		return apperrors.ErrDatabase
	}
	return nil
}

// Delete soft deletes an address
func (r *AddressRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.Address{}).Error; err != nil {
		return apperrors.ErrDatabase
	}
	return nil
}

// SetDefaultAddress sets an address as default and unsets all other addresses for the user
func (r *AddressRepository) SetDefaultAddress(ctx context.Context, userID, addressID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Unset all default addresses for the user
		if err := tx.Model(&domain.Address{}).
			Where("user_id = ? AND is_default = ?", userID, true).
			Update("is_default", false).Error; err != nil {
			return apperrors.ErrDatabase
		}

		// Set the new default address
		if err := tx.Model(&domain.Address{}).
			Where("id = ? AND user_id = ?", addressID, userID).
			Update("is_default", true).Error; err != nil {
			return apperrors.ErrDatabase
		}

		return nil
	})
}

// GetDefaultAddress retrieves the default address for a user
func (r *AddressRepository) GetDefaultAddress(ctx context.Context, userID uuid.UUID) (*domain.Address, error) {
	var address domain.Address
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_default = ?", userID, true).
		First(&address).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrAddressNotFound
		}
		return nil, apperrors.ErrDatabase
	}
	return &address, nil
}

// Exists checks if an address exists
func (r *AddressRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&domain.Address{}).
		Where("id = ?", id).
		Count(&count).Error; err != nil {
		return false, apperrors.ErrDatabase
	}
	return count > 0, nil
}

// BelongsToUser checks if an address belongs to a user
func (r *AddressRepository) BelongsToUser(ctx context.Context, addressID, userID uuid.UUID) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&domain.Address{}).
		Where("id = ? AND user_id = ?", addressID, userID).
		Count(&count).Error; err != nil {
		return false, apperrors.ErrDatabase
	}
	return count > 0, nil
}
