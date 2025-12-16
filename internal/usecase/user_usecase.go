package usecase

import (
	"context"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/oqart/backend/internal/delivery/http/dto"
	"github.com/oqart/backend/internal/domain"
	"github.com/oqart/backend/internal/repository/postgres"
	"github.com/oqart/backend/pkg/cache"
	apperrors "github.com/oqart/backend/pkg/errors"
	"github.com/oqart/backend/pkg/logger"
)

type UserUseCase interface {
	GetProfile(ctx context.Context, userID uuid.UUID) (*dto.UserResponse, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, req *dto.UpdateProfileRequest) (*dto.UserResponse, error)
	ChangePassword(ctx context.Context, userID uuid.UUID, req *dto.ChangePasswordRequest) error
	UpdateEmail(ctx context.Context, userID uuid.UUID, req *dto.UpdateEmailRequest) error
	UpdatePhone(ctx context.Context, userID uuid.UUID, req *dto.UpdatePhoneRequest) error
	DeleteAccount(ctx context.Context, userID uuid.UUID) error
	GetPreferences(ctx context.Context, userID uuid.UUID) (*dto.PreferencesResponse, error)
	UpdatePreferences(ctx context.Context, userID uuid.UUID, req *dto.UpdatePreferencesRequest) (*dto.PreferencesResponse, error)

	// Address management
	GetAddresses(ctx context.Context, userID uuid.UUID) ([]*dto.AddressResponse, error)
	CreateAddress(ctx context.Context, userID uuid.UUID, req *dto.CreateAddressRequest) (*dto.AddressResponse, error)
	UpdateAddress(ctx context.Context, userID, addressID uuid.UUID, req *dto.UpdateAddressRequest) (*dto.AddressResponse, error)
	DeleteAddress(ctx context.Context, userID, addressID uuid.UUID) error
	SetDefaultAddress(ctx context.Context, userID, addressID uuid.UUID) error
}

type userUseCase struct {
	userRepo    *postgres.UserRepository
	addressRepo *postgres.AddressRepository
	cache       *cache.Cache
	logger      *logger.Logger
}

func NewUserUseCase(
	userRepo *postgres.UserRepository,
	addressRepo *postgres.AddressRepository,
	cache *cache.Cache,
	logger *logger.Logger,
) UserUseCase {
	return &userUseCase{
		userRepo:    userRepo,
		addressRepo: addressRepo,
		cache:       cache,
		logger:      logger,
	}
}

// GetProfile retrieves the user's profile
func (uc *userUseCase) GetProfile(ctx context.Context, userID uuid.UUID) (*dto.UserResponse, error) {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		uc.logger.Error("Failed to get user profile", "error", err, "user_id", userID)
		return nil, err
	}

	return uc.mapUserToResponse(user), nil
}

// UpdateProfile updates the user's profile
func (uc *userUseCase) UpdateProfile(ctx context.Context, userID uuid.UUID, req *dto.UpdateProfileRequest) (*dto.UserResponse, error) {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		uc.logger.Error("Failed to get user", "error", err, "user_id", userID)
		return nil, err
	}

	// Update fields
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	if req.Phone != nil {
		user.Phone = req.Phone
	}
	if req.Avatar != nil {
		user.Avatar = req.Avatar
	}

	if err := uc.userRepo.Update(ctx, user); err != nil {
		uc.logger.Error("Failed to update user profile", "error", err, "user_id", userID)
		return nil, err
	}

	uc.logger.Info("User profile updated", "user_id", userID)
	return uc.mapUserToResponse(user), nil
}

// ChangePassword changes the user's password
func (uc *userUseCase) ChangePassword(ctx context.Context, userID uuid.UUID, req *dto.ChangePasswordRequest) error {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		uc.logger.Error("Failed to get user", "error", err, "user_id", userID)
		return err
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword)); err != nil {
		return apperrors.ErrInvalidCredentials
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		uc.logger.Error("Failed to hash password", "error", err)
		return apperrors.ErrInternal
	}

	user.PasswordHash = string(hashedPassword)

	if err := uc.userRepo.Update(ctx, user); err != nil {
		uc.logger.Error("Failed to update password", "error", err, "user_id", userID)
		return err
	}

	uc.logger.Info("Password changed successfully", "user_id", userID)
	return nil
}

// UpdateEmail updates the user's email
func (uc *userUseCase) UpdateEmail(ctx context.Context, userID uuid.UUID, req *dto.UpdateEmailRequest) error {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		uc.logger.Error("Failed to get user", "error", err, "user_id", userID)
		return err
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return apperrors.ErrInvalidCredentials
	}

	// Check if email already exists
	existingUser, err := uc.userRepo.GetByEmail(ctx, req.NewEmail)
	if err != nil && err != apperrors.ErrUserNotFound {
		return err
	}
	if existingUser != nil {
		return apperrors.ErrEmailExists
	}

	user.Email = req.NewEmail
	user.EmailVerified = false // Reset email verification

	if err := uc.userRepo.Update(ctx, user); err != nil {
		uc.logger.Error("Failed to update email", "error", err, "user_id", userID)
		return err
	}

	uc.logger.Info("Email updated successfully", "user_id", userID, "new_email", req.NewEmail)
	return nil
}

// UpdatePhone updates the user's phone number
func (uc *userUseCase) UpdatePhone(ctx context.Context, userID uuid.UUID, req *dto.UpdatePhoneRequest) error {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		uc.logger.Error("Failed to get user", "error", err, "user_id", userID)
		return err
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return apperrors.ErrInvalidCredentials
	}

	// Check if phone already exists
	existingUser, err := uc.userRepo.GetByPhone(ctx, req.NewPhone)
	if err != nil && err != apperrors.ErrUserNotFound {
		return err
	}
	if existingUser != nil {
		return apperrors.ErrPhoneExists
	}

	user.Phone = &req.NewPhone
	user.PhoneVerified = false // Reset phone verification

	if err := uc.userRepo.Update(ctx, user); err != nil {
		uc.logger.Error("Failed to update phone", "error", err, "user_id", userID)
		return err
	}

	uc.logger.Info("Phone updated successfully", "user_id", userID)
	return nil
}

// DeleteAccount soft deletes the user's account
func (uc *userUseCase) DeleteAccount(ctx context.Context, userID uuid.UUID) error {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		uc.logger.Error("Failed to get user", "error", err, "user_id", userID)
		return err
	}

	// Update status to inactive
	user.Status = domain.UserStatusInactive

	if err := uc.userRepo.Update(ctx, user); err != nil {
		uc.logger.Error("Failed to delete account", "error", err, "user_id", userID)
		return err
	}

	uc.logger.Info("Account deleted successfully", "user_id", userID)
	return nil
}

// GetPreferences retrieves the user's preferences
func (uc *userUseCase) GetPreferences(ctx context.Context, userID uuid.UUID) (*dto.PreferencesResponse, error) {
	prefs, err := uc.userRepo.GetPreferences(ctx, userID)
	if err != nil {
		uc.logger.Error("Failed to get preferences", "error", err, "user_id", userID)
		return nil, err
	}

	return uc.mapPreferencesToResponse(prefs), nil
}

// UpdatePreferences updates the user's preferences
func (uc *userUseCase) UpdatePreferences(ctx context.Context, userID uuid.UUID, req *dto.UpdatePreferencesRequest) (*dto.PreferencesResponse, error) {
	prefs, err := uc.userRepo.GetPreferences(ctx, userID)
	if err != nil {
		uc.logger.Error("Failed to get preferences", "error", err, "user_id", userID)
		return nil, err
	}

	// Update only provided fields
	if req.Language != nil {
		prefs.Language = *req.Language
	}
	if req.Currency != nil {
		prefs.Currency = *req.Currency
	}
	if req.NotificationEmail != nil {
		prefs.NotificationEmail = *req.NotificationEmail
	}
	if req.NotificationSMS != nil {
		prefs.NotificationSMS = *req.NotificationSMS
	}
	if req.NotificationPush != nil {
		prefs.NotificationPush = *req.NotificationPush
	}
	if req.MarketingEmails != nil {
		prefs.MarketingEmails = *req.MarketingEmails
	}
	if req.OrderUpdates != nil {
		prefs.OrderUpdates = *req.OrderUpdates
	}
	if req.NewsletterSubscribed != nil {
		prefs.NewsletterSubscribed = *req.NewsletterSubscribed
	}
	if req.TwoFactorEnabled != nil {
		prefs.TwoFactorEnabled = *req.TwoFactorEnabled
	}

	if err := uc.userRepo.UpdatePreferences(ctx, prefs); err != nil {
		uc.logger.Error("Failed to update preferences", "error", err, "user_id", userID)
		return nil, err
	}

	uc.logger.Info("Preferences updated successfully", "user_id", userID)
	return uc.mapPreferencesToResponse(prefs), nil
}

// GetAddresses retrieves all addresses for the user
func (uc *userUseCase) GetAddresses(ctx context.Context, userID uuid.UUID) ([]*dto.AddressResponse, error) {
	addresses, err := uc.addressRepo.GetByUserID(ctx, userID)
	if err != nil {
		uc.logger.Error("Failed to get addresses", "error", err, "user_id", userID)
		return nil, err
	}

	responses := make([]*dto.AddressResponse, len(addresses))
	for i, addr := range addresses {
		responses[i] = uc.mapAddressToResponse(addr)
	}

	return responses, nil
}

// CreateAddress creates a new address for the user
func (uc *userUseCase) CreateAddress(ctx context.Context, userID uuid.UUID, req *dto.CreateAddressRequest) (*dto.AddressResponse, error) {
	address := &domain.Address{
		ID:           uuid.New(),
		UserID:       userID,
		FullName:     req.FullName,
		Phone:        req.Phone,
		AddressLine1: req.AddressLine1,
		AddressLine2: req.AddressLine2,
		City:         req.City,
		State:        req.State,
		Pincode:      req.Pincode,
		Country:      req.Country,
		AddressType:  req.AddressType,
		IsDefault:    req.IsDefault,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
	}

	// If this is marked as default, unset all other defaults first
	if req.IsDefault {
		if err := uc.addressRepo.SetDefaultAddress(ctx, userID, address.ID); err != nil {
			uc.logger.Error("Failed to set default address", "error", err, "user_id", userID)
			return nil, err
		}
	}

	if err := uc.addressRepo.Create(ctx, address); err != nil {
		uc.logger.Error("Failed to create address", "error", err, "user_id", userID)
		return nil, err
	}

	uc.logger.Info("Address created successfully", "user_id", userID, "address_id", address.ID)
	return uc.mapAddressToResponse(address), nil
}

// UpdateAddress updates an existing address
func (uc *userUseCase) UpdateAddress(ctx context.Context, userID, addressID uuid.UUID, req *dto.UpdateAddressRequest) (*dto.AddressResponse, error) {
	// Check if address belongs to user
	belongs, err := uc.addressRepo.BelongsToUser(ctx, addressID, userID)
	if err != nil {
		return nil, err
	}
	if !belongs {
		return nil, apperrors.ErrPermissionDenied
	}

	address, err := uc.addressRepo.GetByID(ctx, addressID)
	if err != nil {
		uc.logger.Error("Failed to get address", "error", err, "address_id", addressID)
		return nil, err
	}

	// Update only provided fields
	if req.FullName != nil {
		address.FullName = *req.FullName
	}
	if req.Phone != nil {
		address.Phone = *req.Phone
	}
	if req.AddressLine1 != nil {
		address.AddressLine1 = *req.AddressLine1
	}
	if req.AddressLine2 != nil {
		address.AddressLine2 = req.AddressLine2
	}
	if req.City != nil {
		address.City = *req.City
	}
	if req.State != nil {
		address.State = *req.State
	}
	if req.Pincode != nil {
		address.Pincode = *req.Pincode
	}
	if req.Country != nil {
		address.Country = *req.Country
	}
	if req.AddressType != nil {
		address.AddressType = *req.AddressType
	}
	if req.Latitude != nil {
		address.Latitude = req.Latitude
	}
	if req.Longitude != nil {
		address.Longitude = req.Longitude
	}

	if err := uc.addressRepo.Update(ctx, address); err != nil {
		uc.logger.Error("Failed to update address", "error", err, "address_id", addressID)
		return nil, err
	}

	uc.logger.Info("Address updated successfully", "user_id", userID, "address_id", addressID)
	return uc.mapAddressToResponse(address), nil
}

// DeleteAddress deletes an address
func (uc *userUseCase) DeleteAddress(ctx context.Context, userID, addressID uuid.UUID) error {
	// Check if address belongs to user
	belongs, err := uc.addressRepo.BelongsToUser(ctx, addressID, userID)
	if err != nil {
		return err
	}
	if !belongs {
		return apperrors.ErrPermissionDenied
	}

	if err := uc.addressRepo.Delete(ctx, addressID); err != nil {
		uc.logger.Error("Failed to delete address", "error", err, "address_id", addressID)
		return err
	}

	uc.logger.Info("Address deleted successfully", "user_id", userID, "address_id", addressID)
	return nil
}

// SetDefaultAddress sets an address as the default
func (uc *userUseCase) SetDefaultAddress(ctx context.Context, userID, addressID uuid.UUID) error {
	// Check if address belongs to user
	belongs, err := uc.addressRepo.BelongsToUser(ctx, addressID, userID)
	if err != nil {
		return err
	}
	if !belongs {
		return apperrors.ErrPermissionDenied
	}

	if err := uc.addressRepo.SetDefaultAddress(ctx, userID, addressID); err != nil {
		uc.logger.Error("Failed to set default address", "error", err, "address_id", addressID)
		return err
	}

	uc.logger.Info("Default address set successfully", "user_id", userID, "address_id", addressID)
	return nil
}

// Helper methods

func (uc *userUseCase) mapUserToResponse(user *domain.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:            user.ID.String(),
		Email:         user.Email,
		Phone:         user.Phone,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Role:          string(user.Role),
		Status:        string(user.Status),
		EmailVerified: user.EmailVerified,
		PhoneVerified: user.PhoneVerified,
		Avatar:        user.Avatar,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
		LastLoginAt:   user.LastLoginAt,
	}
}

func (uc *userUseCase) mapPreferencesToResponse(prefs *domain.UserPreferences) *dto.PreferencesResponse {
	return &dto.PreferencesResponse{
		UserID:               prefs.UserID.String(),
		Language:             prefs.Language,
		Currency:             prefs.Currency,
		NotificationEmail:    prefs.NotificationEmail,
		NotificationSMS:      prefs.NotificationSMS,
		NotificationPush:     prefs.NotificationPush,
		MarketingEmails:      prefs.MarketingEmails,
		OrderUpdates:         prefs.OrderUpdates,
		NewsletterSubscribed: prefs.NewsletterSubscribed,
		TwoFactorEnabled:     prefs.TwoFactorEnabled,
		CreatedAt:            prefs.CreatedAt,
		UpdatedAt:            prefs.UpdatedAt,
	}
}

func (uc *userUseCase) mapAddressToResponse(addr *domain.Address) *dto.AddressResponse {
	return &dto.AddressResponse{
		ID:           addr.ID.String(),
		UserID:       addr.UserID.String(),
		FullName:     addr.FullName,
		Phone:        addr.Phone,
		AddressLine1: addr.AddressLine1,
		AddressLine2: addr.AddressLine2,
		City:         addr.City,
		State:        addr.State,
		Pincode:      addr.Pincode,
		Country:      addr.Country,
		AddressType:  addr.AddressType,
		IsDefault:    addr.IsDefault,
		Latitude:     addr.Latitude,
		Longitude:    addr.Longitude,
		CreatedAt:    addr.CreatedAt,
		UpdatedAt:    addr.UpdatedAt,
	}
}
