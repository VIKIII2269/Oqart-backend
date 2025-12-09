package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/oqart/backend/config"
	"github.com/oqart/backend/internal/delivery/http/dto"
	"github.com/oqart/backend/internal/domain"
	"github.com/oqart/backend/internal/repository/postgres"
	"github.com/oqart/backend/pkg/cache"
	apperrors "github.com/oqart/backend/pkg/errors"
	"github.com/oqart/backend/pkg/jwt"
	"github.com/oqart/backend/pkg/logger"
	"github.com/oqart/backend/pkg/utils"
)

type AuthUseCase interface {
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error)
	SendPhoneOTP(ctx context.Context, req *dto.PhoneLoginRequest) (*dto.OTPResponse, error)
	VerifyPhoneOTP(ctx context.Context, req *dto.VerifyOTPRequest) (*dto.AuthResponse, error)
	RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.AuthResponse, error)
	ForgotPassword(ctx context.Context, req *dto.ForgotPasswordRequest) (*dto.MessageResponse, error)
	ResetPassword(ctx context.Context, req *dto.ResetPasswordRequest) (*dto.MessageResponse, error)
	VerifyEmail(ctx context.Context, token string) (*dto.MessageResponse, error)
	Logout(ctx context.Context, userID uuid.UUID) error
}

type authUseCase struct {
	userRepo          postgres.UserRepository
	sessionRepo       postgres.SessionRepository
	otpRepo           postgres.OTPRepository
	passwordResetRepo postgres.PasswordResetRepository
	jwtManager        *jwt.TokenManager
	cache             *cache.Cache
	config            *config.Config
	logger            *logger.Logger
	// TODO: Add email and SMS services
}

func NewAuthUseCase(
	userRepo postgres.UserRepository,
	sessionRepo postgres.SessionRepository,
	otpRepo postgres.OTPRepository,
	passwordResetRepo postgres.PasswordResetRepository,
	jwtManager *jwt.TokenManager,
	cache *cache.Cache,
	config *config.Config,
	logger *logger.Logger,
) AuthUseCase {
	return &authUseCase{
		userRepo:          userRepo,
		sessionRepo:       sessionRepo,
		otpRepo:           otpRepo,
		passwordResetRepo: passwordResetRepo,
		jwtManager:        jwtManager,
		cache:             cache,
		config:            config,
		logger:            logger,
	}
}

func (uc *authUseCase) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	// Check if email already exists
	existingUser, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, apperrors.ErrEmailExists
	}

	// Hash password
	passwordHash, err := utils.HashPassword(req.Password, uc.config.Security.BcryptCost)
	if err != nil {
		uc.logger.Error("Failed to hash password", err)
		return nil, apperrors.ErrInternal
	}

	// Create user
	user := &domain.User{
		Email:        &req.Email,
		PasswordHash: &passwordHash,
		Role:         domain.RoleCustomer,
		Status:       domain.StatusActive,
		AuthProvider: domain.AuthProviderEmail,
		FirstName:    &req.FirstName,
		LastName:     &req.LastName,
	}

	if req.Phone != "" {
		sanitizedPhone := utils.SanitizePhone(req.Phone)
		user.Phone = &sanitizedPhone
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		uc.logger.Error("Failed to create user", err)
		return nil, apperrors.ErrInternal
	}

	// Generate tokens
	tokens, err := uc.jwtManager.GenerateTokenPair(user.ID.String(), *user.Email, string(user.Role))
	if err != nil {
		uc.logger.Error("Failed to generate tokens", err)
		return nil, apperrors.ErrInternal
	}

	// Create session
	session := &domain.Session{
		UserID:           user.ID,
		AccessToken:      tokens.AccessToken,
		RefreshToken:     tokens.RefreshToken,
		ExpiresAt:        tokens.ExpiresAt,
		RefreshExpiresAt: time.Now().Add(uc.config.JWT.RefreshTokenExpiry),
	}

	if err := uc.sessionRepo.Create(ctx, session); err != nil {
		uc.logger.Error("Failed to create session", err)
		return nil, apperrors.ErrInternal
	}

	// TODO: Send verification email

	uc.logger.Info("User registered successfully", "user_id", user.ID.String(), "email", req.Email)

	return &dto.AuthResponse{
		User:         dto.ToUserResponse(user),
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    tokens.ExpiresAt,
	}, nil
}

func (uc *authUseCase) Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error) {
	// Get user by email
	user, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, apperrors.ErrInvalidCredentials
	}

	// Check if user is active
	if !user.IsActive() {
		return nil, apperrors.New("ACCOUNT_INACTIVE", "Your account is inactive", 403)
	}

	// Verify password
	if user.PasswordHash == nil || !utils.CheckPassword(req.Password, *user.PasswordHash) {
		return nil, apperrors.ErrInvalidCredentials
	}

	// Generate tokens
	tokens, err := uc.jwtManager.GenerateTokenPair(user.ID.String(), *user.Email, string(user.Role))
	if err != nil {
		uc.logger.Error("Failed to generate tokens", err)
		return nil, apperrors.ErrInternal
	}

	// Create session
	session := &domain.Session{
		UserID:           user.ID,
		AccessToken:      tokens.AccessToken,
		RefreshToken:     tokens.RefreshToken,
		ExpiresAt:        tokens.ExpiresAt,
		RefreshExpiresAt: time.Now().Add(uc.config.JWT.RefreshTokenExpiry),
	}

	if err := uc.sessionRepo.Create(ctx, session); err != nil {
		uc.logger.Error("Failed to create session", err)
		return nil, apperrors.ErrInternal
	}

	// Update last login
	if err := uc.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		uc.logger.Warn("Failed to update last login", err)
	}

	uc.logger.Info("User logged in successfully", "user_id", user.ID.String(), "email", req.Email)

	return &dto.AuthResponse{
		User:         dto.ToUserResponse(user),
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    tokens.ExpiresAt,
	}, nil
}

func (uc *authUseCase) SendPhoneOTP(ctx context.Context, req *dto.PhoneLoginRequest) (*dto.OTPResponse, error) {
	sanitizedPhone := utils.SanitizePhone(req.Phone)

	// Check rate limiting
	cacheKey := fmt.Sprintf("otp:rate_limit:%s", sanitizedPhone)
	if uc.cache != nil {
		count, err := uc.cache.Get(ctx, cacheKey, new(int))
		if err == nil && count != nil {
			if *count.(*int) >= 5 {
				return nil, apperrors.New("RATE_LIMIT", "Too many OTP requests. Please try again later.", 429)
			}
		}
	}

	// Generate OTP
	otpCode, err := utils.Generate6DigitOTP()
	if err != nil {
		uc.logger.Error("Failed to generate OTP", err)
		return nil, apperrors.ErrInternal
	}

	expiresAt := time.Now().Add(uc.config.Security.OTPExpiry)

	// Save OTP
	otp := &domain.OTP{
		Phone:     sanitizedPhone,
		Code:      otpCode,
		ExpiresAt: expiresAt,
	}

	if err := uc.otpRepo.Create(ctx, otp); err != nil {
		uc.logger.Error("Failed to create OTP", err)
		return nil, apperrors.ErrInternal
	}

	// TODO: Send OTP via SMS
	uc.logger.Info("OTP generated", "phone", sanitizedPhone, "code", otpCode) // Remove in production

	// Update rate limit
	if uc.cache != nil {
		currentCount := 0
		uc.cache.Get(ctx, cacheKey, &currentCount)
		uc.cache.Set(ctx, cacheKey, currentCount+1, 1*time.Hour)
	}

	return &dto.OTPResponse{
		Message:   "OTP sent successfully to your phone",
		ExpiresAt: expiresAt,
	}, nil
}

func (uc *authUseCase) VerifyPhoneOTP(ctx context.Context, req *dto.VerifyOTPRequest) (*dto.AuthResponse, error) {
	sanitizedPhone := utils.SanitizePhone(req.Phone)

	// Get latest OTP
	otp, err := uc.otpRepo.GetLatestByPhone(ctx, sanitizedPhone)
	if err != nil {
		return nil, apperrors.New("INVALID_OTP", "Invalid or expired OTP", 400)
	}

	// Check attempts
	if otp.Attempts >= 3 {
		return nil, apperrors.New("TOO_MANY_ATTEMPTS", "Too many failed attempts. Please request a new OTP.", 400)
	}

	// Verify OTP
	if otp.Code != req.Code {
		uc.otpRepo.IncrementAttempts(ctx, otp.ID)
		return nil, apperrors.New("INVALID_OTP", "Invalid OTP", 400)
	}

	// Mark OTP as verified
	if err := uc.otpRepo.MarkAsVerified(ctx, otp.ID); err != nil {
		uc.logger.Error("Failed to mark OTP as verified", err)
	}

	// Get or create user
	user, err := uc.userRepo.GetByPhone(ctx, sanitizedPhone)
	if err != nil {
		// Create new user
		user = &domain.User{
			Phone:        &sanitizedPhone,
			Role:         domain.RoleCustomer,
			Status:       domain.StatusActive,
			AuthProvider: domain.AuthProviderPhone,
			PhoneVerified: true,
		}

		if err := uc.userRepo.Create(ctx, user); err != nil {
			uc.logger.Error("Failed to create user", err)
			return nil, apperrors.ErrInternal
		}
	} else {
		// Verify phone if not already verified
		if !user.PhoneVerified {
			if err := uc.userRepo.VerifyPhone(ctx, user.ID); err != nil {
				uc.logger.Warn("Failed to verify phone", err)
			}
		}
	}

	// Generate tokens
	email := ""
	if user.Email != nil {
		email = *user.Email
	}
	tokens, err := uc.jwtManager.GenerateTokenPair(user.ID.String(), email, string(user.Role))
	if err != nil {
		uc.logger.Error("Failed to generate tokens", err)
		return nil, apperrors.ErrInternal
	}

	// Create session
	session := &domain.Session{
		UserID:           user.ID,
		AccessToken:      tokens.AccessToken,
		RefreshToken:     tokens.RefreshToken,
		ExpiresAt:        tokens.ExpiresAt,
		RefreshExpiresAt: time.Now().Add(uc.config.JWT.RefreshTokenExpiry),
	}

	if err := uc.sessionRepo.Create(ctx, session); err != nil {
		uc.logger.Error("Failed to create session", err)
		return nil, apperrors.ErrInternal
	}

	// Update last login
	if err := uc.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		uc.logger.Warn("Failed to update last login", err)
	}

	uc.logger.Info("User logged in with OTP successfully", "user_id", user.ID.String(), "phone", sanitizedPhone)

	return &dto.AuthResponse{
		User:         dto.ToUserResponse(user),
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    tokens.ExpiresAt,
	}, nil
}

func (uc *authUseCase) RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.AuthResponse, error) {
	// Get session by refresh token
	session, err := uc.sessionRepo.GetByRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, apperrors.New("INVALID_TOKEN", "Invalid refresh token", 401)
	}

	// Check if refresh token is expired
	if time.Now().After(session.RefreshExpiresAt) {
		uc.sessionRepo.Delete(ctx, session.ID)
		return nil, apperrors.New("TOKEN_EXPIRED", "Refresh token expired", 401)
	}

	// Get user
	user := &session.User

	// Generate new tokens
	email := ""
	if user.Email != nil {
		email = *user.Email
	}
	tokens, err := uc.jwtManager.GenerateTokenPair(user.ID.String(), email, string(user.Role))
	if err != nil {
		uc.logger.Error("Failed to generate tokens", err)
		return nil, apperrors.ErrInternal
	}

	// Update session
	session.AccessToken = tokens.AccessToken
	session.RefreshToken = tokens.RefreshToken
	session.ExpiresAt = tokens.ExpiresAt
	session.RefreshExpiresAt = time.Now().Add(uc.config.JWT.RefreshTokenExpiry)

	if err := uc.sessionRepo.Create(ctx, session); err != nil {
		uc.logger.Error("Failed to update session", err)
		return nil, apperrors.ErrInternal
	}

	return &dto.AuthResponse{
		User:         dto.ToUserResponse(user),
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    tokens.ExpiresAt,
	}, nil
}

func (uc *authUseCase) ForgotPassword(ctx context.Context, req *dto.ForgotPasswordRequest) (*dto.MessageResponse, error) {
	// Get user
	user, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		// Don't reveal if email exists
		return &dto.MessageResponse{Message: "If the email exists, a password reset link has been sent"}, nil
	}

	// Generate reset token
	token, err := utils.GenerateSecureToken()
	if err != nil {
		uc.logger.Error("Failed to generate reset token", err)
		return nil, apperrors.ErrInternal
	}

	// Save password reset
	reset := &domain.PasswordReset{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(uc.config.Security.PasswordResetExpiry),
	}

	if err := uc.passwordResetRepo.Create(ctx, reset); err != nil {
		uc.logger.Error("Failed to create password reset", err)
		return nil, apperrors.ErrInternal
	}

	// TODO: Send password reset email with token

	uc.logger.Info("Password reset requested", "user_id", user.ID.String(), "email", req.Email)

	return &dto.MessageResponse{Message: "If the email exists, a password reset link has been sent"}, nil
}

func (uc *authUseCase) ResetPassword(ctx context.Context, req *dto.ResetPasswordRequest) (*dto.MessageResponse, error) {
	// Get password reset by token
	reset, err := uc.passwordResetRepo.GetByToken(ctx, req.Token)
	if err != nil {
		return nil, apperrors.New("INVALID_TOKEN", "Invalid or expired reset token", 400)
	}

	// Hash new password
	passwordHash, err := utils.HashPassword(req.Password, uc.config.Security.BcryptCost)
	if err != nil {
		uc.logger.Error("Failed to hash password", err)
		return nil, apperrors.ErrInternal
	}

	// Update password
	if err := uc.userRepo.ChangePassword(ctx, reset.UserID, passwordHash); err != nil {
		uc.logger.Error("Failed to change password", err)
		return nil, apperrors.ErrInternal
	}

	// Mark reset as used
	if err := uc.passwordResetRepo.MarkAsUsed(ctx, reset.ID); err != nil {
		uc.logger.Warn("Failed to mark reset as used", err)
	}

	// Delete all user sessions (force re-login)
	if err := uc.sessionRepo.DeleteByUserID(ctx, reset.UserID); err != nil {
		uc.logger.Warn("Failed to delete user sessions", err)
	}

	uc.logger.Info("Password reset successfully", "user_id", reset.UserID.String())

	return &dto.MessageResponse{Message: "Password reset successfully"}, nil
}

func (uc *authUseCase) VerifyEmail(ctx context.Context, token string) (*dto.MessageResponse, error) {
	// TODO: Implement email verification
	return &dto.MessageResponse{Message: "Email verified successfully"}, nil
}

func (uc *authUseCase) Logout(ctx context.Context, userID uuid.UUID) error {
	// Delete all user sessions
	if err := uc.sessionRepo.DeleteByUserID(ctx, userID); err != nil {
		uc.logger.Error("Failed to delete user sessions", err)
		return apperrors.ErrInternal
	}

	uc.logger.Info("User logged out successfully", "user_id", userID.String())
	return nil
}
