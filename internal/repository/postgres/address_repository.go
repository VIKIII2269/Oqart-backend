package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/oqart/backend/internal/domain"
	"github.com/oqart/backend/pkg/errors"
	"gorm.io/gorm"
)

type AddressRepository interface {
	Create(ctx context.Context, address *domain.Address) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Address, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Address, error)
	GetDefaultByUserID(ctx context.Context, userID uuid.UUID) (*domain.Address, error)
	Update(ctx context.Context, address *domain.Address) error
	Delete(ctx context.Context, id uuid.UUID) error
	SetAsDefault(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

type addressRepository struct {
	db *gorm.DB
}

func NewAddressRepository(db *gorm.DB) AddressRepository {
	return &addressRepository{db: db}
}

func (r *addressRepository) Create(ctx context.Context, address *domain.Address) error {
	// If this is set as default, unset others
	if address.IsDefault {
		if err := r.db.WithContext(ctx).Model(&domain.Address{}).
			Where("user_id = ? AND is_default = true", address.UserID).
			Update("is_default", false).Error; err != nil {
			return fmt.Errorf("failed to unset default addresses: %w", err)
		}
	}

	if err := r.db.WithContext(ctx).Create(address).Error; err != nil {
		return fmt.Errorf("failed to create address: %w", err)
	}
	return nil
}

func (r *addressRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Address, error) {
	var address domain.Address
	if err := r.db.WithContext(ctx).First(&address, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get address: %w", err)
	}
	return &address, nil
}

func (r *addressRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Address, error) {
	var addresses []*domain.Address
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("is_default DESC, created_at DESC").
		Find(&addresses).Error; err != nil {
		return nil, fmt.Errorf("failed to get addresses: %w", err)
	}
	return addresses, nil
}

func (r *addressRepository) GetDefaultByUserID(ctx context.Context, userID uuid.UUID) (*domain.Address, error) {
	var address domain.Address
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_default = true", userID).
		First(&address).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get default address: %w", err)
	}
	return &address, nil
}

func (r *addressRepository) Update(ctx context.Context, address *domain.Address) error {
	// If this is being set as default, unset others
	if address.IsDefault {
		if err := r.db.WithContext(ctx).Model(&domain.Address{}).
			Where("user_id = ? AND id != ? AND is_default = true", address.UserID, address.ID).
			Update("is_default", false).Error; err != nil {
			return fmt.Errorf("failed to unset default addresses: %w", err)
		}
	}

	if err := r.db.WithContext(ctx).Save(address).Error; err != nil {
		return fmt.Errorf("failed to update address: %w", err)
	}
	return nil
}

func (r *addressRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&domain.Address{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("failed to delete address: %w", err)
	}
	return nil
}

func (r *addressRepository) SetAsDefault(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Unset all default addresses for user
		if err := tx.Model(&domain.Address{}).
			Where("user_id = ? AND is_default = true", userID).
			Update("is_default", false).Error; err != nil {
			return fmt.Errorf("failed to unset default addresses: %w", err)
		}

		// Set this address as default
		if err := tx.Model(&domain.Address{}).
			Where("id = ? AND user_id = ?", id, userID).
			Update("is_default", true).Error; err != nil {
			return fmt.Errorf("failed to set default address: %w", err)
		}

		return nil
	})
}

// User Preferences Repository
type UserPreferencesRepository interface {
	Create(ctx context.Context, prefs *domain.UserPreferences) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.UserPreferences, error)
	Update(ctx context.Context, prefs *domain.UserPreferences) error
	Upsert(ctx context.Context, prefs *domain.UserPreferences) error
}

type userPreferencesRepository struct {
	db *gorm.DB
}

func NewUserPreferencesRepository(db *gorm.DB) UserPreferencesRepository {
	return &userPreferencesRepository{db: db}
}

func (r *userPreferencesRepository) Create(ctx context.Context, prefs *domain.UserPreferences) error {
	if err := r.db.WithContext(ctx).Create(prefs).Error; err != nil {
		return fmt.Errorf("failed to create user preferences: %w", err)
	}
	return nil
}

func (r *userPreferencesRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.UserPreferences, error) {
	var prefs domain.UserPreferences
	if err := r.db.WithContext(ctx).First(&prefs, "user_id = ?", userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create default preferences if not found
			defaultPrefs := &domain.UserPreferences{
				UserID:             userID,
				Language:           "en",
				Currency:           "INR",
				Timezone:           "Asia/Kolkata",
				EmailNotifications: true,
				SMSNotifications:   true,
				PushNotifications:  true,
				MarketingEmails:    true,
			}
			if err := r.Create(ctx, defaultPrefs); err != nil {
				return nil, err
			}
			return defaultPrefs, nil
		}
		return nil, fmt.Errorf("failed to get user preferences: %w", err)
	}
	return &prefs, nil
}

func (r *userPreferencesRepository) Update(ctx context.Context, prefs *domain.UserPreferences) error {
	if err := r.db.WithContext(ctx).Save(prefs).Error; err != nil {
		return fmt.Errorf("failed to update user preferences: %w", err)
	}
	return nil
}

func (r *userPreferencesRepository) Upsert(ctx context.Context, prefs *domain.UserPreferences) error {
	// Try to get existing
	existing, err := r.GetByUserID(ctx, prefs.UserID)
	if err != nil {
		return r.Create(ctx, prefs)
	}

	// Update existing
	prefs.ID = existing.ID
	return r.Update(ctx, prefs)
}
