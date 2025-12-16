package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/oqart/backend/config"
	"github.com/oqart/backend/internal/delivery/http/dto"
	"github.com/oqart/backend/internal/domain"
	"github.com/oqart/backend/internal/repository/postgres"
	apperrors "github.com/oqart/backend/pkg/errors"
	"github.com/oqart/backend/pkg/logger"
	"github.com/oqart/backend/pkg/utils"
)

type UserUseCase interface {
	GetProfile(ctx context.Context, userID uuid.UUID) (*dto.UserResponse, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, req *dto.UpdateProfileRequest) (*dto.UserResponse, error)
	ChangePassword(ctx context.Context, userID uuid.UUID, req *dto.ChangePasswordRequest) (*dto.MessageResponse, error)
	UpdateEmail(ctx context.Context, userID uuid.UUID, req *dto.UpdateEmailRequest) (*dto.MessageResponse, error)
	UpdatePhone(ctx context.Context, userID uuid.UUID, req *dto.UpdatePhoneRequest) (*dto.MessageResponse, error)
	DeleteAccount(ctx context.Context, userID uuid.UUID) (*dto.MessageResponse, error)
	GetPreferences(ctx context.Context, userID uuid.UUID) (*dto.UserPreferencesResponse, error)
	UpdatePreferences(ctx context.Context, userID uuid.UUID, req *dto.UpdatePreferencesRequest) (*dto.UserPreferencesResponse, error)

	// Address management
	CreateAddress(ctx context.Context, userID uuid.UUID, req *dto.AddressRequest) (*dto.AddressResponse, error)
	GetAddresses(ctx context.Context, userID uuid.UUID) ([]*dto.AddressResponse, error)
	GetAddress(ctx context.Context, userID uuid.UUID, addressID uuid.UUID) (*dto.AddressResponse, error)
	UpdateAddress(ctx context.Context, userID uuid.UUID, addressID uuid.UUID, req *dto.AddressRequest) (*dto.AddressResponse, error)
	DeleteAddress(ctx context.Context, userID uuid.UUID, addressID uuid.UUID) (*dto.MessageResponse, error)
}

type userUseCase struct {
	userRepo        postgres.UserRepository
	addressRepo     postgres.AddressRepository
	preferencesRepo postgres.UserPreferencesRepository
	sessionRepo     postgres.SessionRepository
	config          *config.Config
	logger          *logger.Logger
}

func NewUserUseCase(
	userRepo postgres.UserRepository,
	addressRepo postgres.AddressRepository,
	preferencesRepo postgres.UserPreferencesRepository,
	sessionRepo postgres.SessionRepository,
	config *config.Config,
	logger *logger.Logger,
) UserUseCase {
	return &userUseCase{
		userRepo:        userRepo,
		addressRepo:     addressRepo,
		preferencesRepo: preferencesRepo,
		sessionRepo:     sessionRepo,
		config:          config,
		logger:          logger,
	}
}

func (uc *userUseCase) GetProfile(ctx context.Context, userID uuid.UUID) (*dto.UserResponse, error) {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return dto.ToUserResponse(user), nil
}

func (uc *userUseCase) UpdateProfile(ctx context.Context, userID uuid.UUID, req *dto.UpdateProfileRequest) (*dto.UserResponse, error) {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.FirstName != nil {
		user.FirstName = req.FirstName
	}
	if req.LastName != nil {
		user.LastName = req.LastName
	}
	if req.DateOfBirth != nil {
		user.DateOfBirth = req.DateOfBirth
	}
	if req.Gender != nil {
		user.Gender = req.Gender
	}

	if err := uc.userRepo.Update(ctx, user); err != nil {
		uc.logger.Error("Failed to update user profile", err, "user_id", userID.String())
		return nil, apperrors.ErrInternal
	}

	uc.logger.Info("User profile updated", "user_id", userID.String())

	return dto.ToUserResponse(user), nil
}

func (uc *userUseCase) ChangePassword(ctx context.Context, userID uuid.UUID, req *dto.ChangePasswordRequest) (*dto.MessageResponse, error) {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Verify current password
	if user.PasswordHash == nil || !utils.CheckPassword(req.CurrentPassword, *user.PasswordHash) {
		return nil, apperrors.New("INVALID_PASSWORD", "Current password is incorrect", 400)
	}

	// Hash new password
	newPasswordHash, err := utils.HashPassword(req.NewPassword, uc.config.Security.BcryptCost)
	if err != nil {
		uc.logger.Error("Failed to hash password", err)
		return nil, apperrors.ErrInternal
	}

	// Update password
	if err := uc.userRepo.ChangePassword(ctx, userID, newPasswordHash); err != nil {
		uc.logger.Error("Failed to change password", err, "user_id", userID.String())
		return nil, apperrors.ErrInternal
	}

	// Delete all sessions to force re-login
	if err := uc.sessionRepo.DeleteByUserID(ctx, userID); err != nil {
		uc.logger.Warn("Failed to delete user sessions", err)
	}

	uc.logger.Info("Password changed successfully", "user_id", userID.String())

	return &dto.MessageResponse{Message: "Password changed successfully. Please login again."}, nil
}

func (uc *userUseCase) UpdateEmail(ctx context.Context, userID uuid.UUID, req *dto.UpdateEmailRequest) (*dto.MessageResponse, error) {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Check if email already exists
	existingUser, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil && existingUser.ID != userID {
		return nil, apperrors.ErrEmailExists
	}

	// Update email
	user.Email = &req.Email
	user.EmailVerified = false // Require re-verification

	if err := uc.userRepo.Update(ctx, user); err != nil {
		uc.logger.Error("Failed to update email", err, "user_id", userID.String())
		return nil, apperrors.ErrInternal
	}

	// TODO: Send email verification

	uc.logger.Info("Email updated", "user_id", userID.String(), "new_email", req.Email)

	return &dto.MessageResponse{Message: "Email updated successfully. Please verify your new email."}, nil
}

func (uc *userUseCase) UpdatePhone(ctx context.Context, userID uuid.UUID, req *dto.UpdatePhoneRequest) (*dto.MessageResponse, error) {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	sanitizedPhone := utils.SanitizePhone(req.Phone)

	// Check if phone already exists
	existingUser, err := uc.userRepo.GetByPhone(ctx, sanitizedPhone)
	if err == nil && existingUser != nil && existingUser.ID != userID {
		return nil, apperrors.ErrPhoneExists
	}

	// Update phone
	user.Phone = &sanitizedPhone
	user.PhoneVerified = false // Require re-verification

	if err := uc.userRepo.Update(ctx, user); err != nil {
		uc.logger.Error("Failed to update phone", err, "user_id", userID.String())
		return nil, apperrors.ErrInternal
	}

	// TODO: Send OTP for verification

	uc.logger.Info("Phone updated", "user_id", userID.String(), "new_phone", sanitizedPhone)

	return &dto.MessageResponse{Message: "Phone number updated successfully. Please verify your new phone number."}, nil
}

func (uc *userUseCase) DeleteAccount(ctx context.Context, userID uuid.UUID) (*dto.MessageResponse, error) {
	// Soft delete by setting status to deleted
	if err := uc.userRepo.Delete(ctx, userID); err != nil {
		uc.logger.Error("Failed to delete account", err, "user_id", userID.String())
		return nil, apperrors.ErrInternal
	}

	// Delete all sessions
	if err := uc.sessionRepo.DeleteByUserID(ctx, userID); err != nil {
		uc.logger.Warn("Failed to delete user sessions", err)
	}

	uc.logger.Info("Account deleted", "user_id", userID.String())

	return &dto.MessageResponse{Message: "Account deleted successfully"}, nil
}

func (uc *userUseCase) GetPreferences(ctx context.Context, userID uuid.UUID) (*dto.UserPreferencesResponse, error) {
	prefs, err := uc.preferencesRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return dto.ToUserPreferencesResponse(prefs), nil
}

func (uc *userUseCase) UpdatePreferences(ctx context.Context, userID uuid.UUID, req *dto.UpdatePreferencesRequest) (*dto.UserPreferencesResponse, error) {
	prefs, err := uc.preferencesRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Language != nil {
		prefs.Language = *req.Language
	}
	if req.Currency != nil {
		prefs.Currency = *req.Currency
	}
	if req.Timezone != nil {
		prefs.Timezone = *req.Timezone
	}
	if req.EmailNotifications != nil {
		prefs.EmailNotifications = *req.EmailNotifications
	}
	if req.SMSNotifications != nil {
		prefs.SMSNotifications = *req.SMSNotifications
	}
	if req.PushNotifications != nil {
		prefs.PushNotifications = *req.PushNotifications
	}
	if req.MarketingEmails != nil {
		prefs.MarketingEmails = *req.MarketingEmails
	}

	if err := uc.preferencesRepo.Update(ctx, prefs); err != nil {
		uc.logger.Error("Failed to update preferences", err, "user_id", userID.String())
		return nil, apperrors.ErrInternal
	}

	uc.logger.Info("User preferences updated", "user_id", userID.String())

	return dto.ToUserPreferencesResponse(prefs), nil
}

// Address management

func (uc *userUseCase) CreateAddress(ctx context.Context, userID uuid.UUID, req *dto.AddressRequest) (*dto.AddressResponse, error) {
	address := &domain.Address{
		UserID:        userID,
		Type:          domain.AddressType(req.Type),
		FullName:      req.FullName,
		Phone:         req.Phone,
		StreetAddress: req.StreetAddress,
		City:          req.City,
		State:         req.State,
		Pincode:       req.Pincode,
		Country:       utils.DefaultString(req.Country, "India"),
		IsDefault:     req.IsDefault,
	}

	if req.Email != "" {
		address.Email = &req.Email
	}
	if req.Landmark != "" {
		address.Landmark = &req.Landmark
	}

	if err := uc.addressRepo.Create(ctx, address); err != nil {
		uc.logger.Error("Failed to create address", err, "user_id", userID.String())
		return nil, apperrors.ErrInternal
	}

	uc.logger.Info("Address created", "user_id", userID.String(), "address_id", address.ID.String())

	return dto.ToAddressResponse(address), nil
}

func (uc *userUseCase) GetAddresses(ctx context.Context, userID uuid.UUID) ([]*dto.AddressResponse, error) {
	addresses, err := uc.addressRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	response := make([]*dto.AddressResponse, len(addresses))
	for i, addr := range addresses {
		response[i] = dto.ToAddressResponse(addr)
	}

	return response, nil
}

func (uc *userUseCase) GetAddress(ctx context.Context, userID uuid.UUID, addressID uuid.UUID) (*dto.AddressResponse, error) {
	address, err := uc.addressRepo.GetByID(ctx, addressID)
	if err != nil {
		return nil, err
	}

	// Verify ownership
	if address.UserID != userID {
		return nil, apperrors.ErrForbidden
	}

	return dto.ToAddressResponse(address), nil
}

func (uc *userUseCase) UpdateAddress(ctx context.Context, userID uuid.UUID, addressID uuid.UUID, req *dto.AddressRequest) (*dto.AddressResponse, error) {
	address, err := uc.addressRepo.GetByID(ctx, addressID)
	if err != nil {
		return nil, err
	}

	// Verify ownership
	if address.UserID != userID {
		return nil, apperrors.ErrForbidden
	}

	// Update fields
	address.Type = domain.AddressType(req.Type)
	address.FullName = req.FullName
	address.Phone = req.Phone
	address.StreetAddress = req.StreetAddress
	address.City = req.City
	address.State = req.State
	address.Pincode = req.Pincode
	address.Country = utils.DefaultString(req.Country, "India")
	address.IsDefault = req.IsDefault

	if req.Email != "" {
		address.Email = &req.Email
	} else {
		address.Email = nil
	}

	if req.Landmark != "" {
		address.Landmark = &req.Landmark
	} else {
		address.Landmark = nil
	}

	if err := uc.addressRepo.Update(ctx, address); err != nil {
		uc.logger.Error("Failed to update address", err, "address_id", addressID.String())
		return nil, apperrors.ErrInternal
	}

	uc.logger.Info("Address updated", "user_id", userID.String(), "address_id", addressID.String())

	return dto.ToAddressResponse(address), nil
}

func (uc *userUseCase) DeleteAddress(ctx context.Context, userID uuid.UUID, addressID uuid.UUID) (*dto.MessageResponse, error) {
	address, err := uc.addressRepo.GetByID(ctx, addressID)
	if err != nil {
		return nil, err
	}

	// Verify ownership
	if address.UserID != userID {
		return nil, apperrors.ErrForbidden
	}

	if err := uc.addressRepo.Delete(ctx, addressID); err != nil {
		uc.logger.Error("Failed to delete address", err, "address_id", addressID.String())
		return nil, apperrors.ErrInternal
	}

	uc.logger.Info("Address deleted", "user_id", userID.String(), "address_id", addressID.String())

	return &dto.MessageResponse{Message: "Address deleted successfully"}, nil
}
