package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/oqart/backend/internal/domain"
	"github.com/oqart/backend/pkg/errors"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByPhone(ctx context.Context, phone string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateLastLogin(ctx context.Context, id uuid.UUID) error
	VerifyEmail(ctx context.Context, id uuid.UUID) error
	VerifyPhone(ctx context.Context, id uuid.UUID) error
	ChangePassword(ctx context.Context, id uuid.UUID, passwordHash string) error
	List(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]*domain.User, int64, error)
	GetPreferences(ctx context.Context, userID uuid.UUID) (*domain.UserPreferences, error)
	UpdatePreferences(ctx context.Context, preferences *domain.UserPreferences) error
	CreatePreferences(ctx context.Context, preferences *domain.UserPreferences) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).First(&user, "email = ?", email).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return &user, nil
}

func (r *userRepository) GetByPhone(ctx context.Context, phone string) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).First(&user, "phone = ?", phone).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by phone: %w", err)
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	// Soft delete - update status instead of actual delete
	if err := r.db.WithContext(ctx).Model(&domain.User{}).Where("id = ?", id).Update("status", domain.StatusDeleted).Error; err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

func (r *userRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	if err := r.db.WithContext(ctx).Model(&domain.User{}).Where("id = ?", id).Update("last_login_at", now).Error; err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}
	return nil
}

func (r *userRepository) VerifyEmail(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Model(&domain.User{}).Where("id = ?", id).Update("email_verified", true).Error; err != nil {
		return fmt.Errorf("failed to verify email: %w", err)
	}
	return nil
}

func (r *userRepository) VerifyPhone(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Model(&domain.User{}).Where("id = ?", id).Update("phone_verified", true).Error; err != nil {
		return fmt.Errorf("failed to verify phone: %w", err)
	}
	return nil
}

func (r *userRepository) ChangePassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	if err := r.db.WithContext(ctx).Model(&domain.User{}).Where("id = ?", id).Update("password_hash", passwordHash).Error; err != nil {
		return fmt.Errorf("failed to change password: %w", err)
	}
	return nil
}

func (r *userRepository) List(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]*domain.User, int64, error) {
	var users []*domain.User
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.User{})

	// Apply filters
	for key, value := range filters {
		query = query.Where(key+" = ?", value)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Get paginated results
	if err := query.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	return users, total, nil
}

func (r *userRepository) GetPreferences(ctx context.Context, userID uuid.UUID) (*domain.UserPreferences, error) {
	var prefs domain.UserPreferences
	if err := r.db.WithContext(ctx).First(&prefs, "user_id = ?", userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create default preferences if none exist
			prefs = domain.UserPreferences{
				ID:                   uuid.New(),
				UserID:               userID,
				Language:             "en",
				Currency:             "INR",
				NotificationEmail:    true,
				NotificationSMS:      true,
				NotificationPush:     true,
				MarketingEmails:      false,
				OrderUpdates:         true,
				NewsletterSubscribed: false,
				TwoFactorEnabled:     false,
			}
			if err := r.db.WithContext(ctx).Create(&prefs).Error; err != nil {
				return nil, fmt.Errorf("failed to create default preferences: %w", err)
			}
			return &prefs, nil
		}
		return nil, fmt.Errorf("failed to get preferences: %w", err)
	}
	return &prefs, nil
}

func (r *userRepository) UpdatePreferences(ctx context.Context, preferences *domain.UserPreferences) error {
	if err := r.db.WithContext(ctx).Save(preferences).Error; err != nil {
		return fmt.Errorf("failed to update preferences: %w", err)
	}
	return nil
}

func (r *userRepository) CreatePreferences(ctx context.Context, preferences *domain.UserPreferences) error {
	if err := r.db.WithContext(ctx).Create(preferences).Error; err != nil {
		return fmt.Errorf("failed to create preferences: %w", err)
	}
	return nil
}

// Session Repository
type SessionRepository interface {
	Create(ctx context.Context, session *domain.Session) error
	GetByToken(ctx context.Context, token string) (*domain.Session, error)
	GetByRefreshToken(ctx context.Context, token string) (*domain.Session, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Session, error)
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
	DeleteExpired(ctx context.Context) error
}

type sessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) SessionRepository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) Create(ctx context.Context, session *domain.Session) error {
	if err := r.db.WithContext(ctx).Create(session).Error; err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	return nil
}

func (r *sessionRepository) GetByToken(ctx context.Context, token string) (*domain.Session, error) {
	var session domain.Session
	if err := r.db.WithContext(ctx).Preload("User").First(&session, "access_token = ?", token).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}
	return &session, nil
}

func (r *sessionRepository) GetByRefreshToken(ctx context.Context, token string) (*domain.Session, error) {
	var session domain.Session
	if err := r.db.WithContext(ctx).Preload("User").First(&session, "refresh_token = ?", token).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get session by refresh token: %w", err)
	}
	return &session, nil
}

func (r *sessionRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Session, error) {
	var sessions []*domain.Session
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&sessions).Error; err != nil {
		return nil, fmt.Errorf("failed to get sessions by user: %w", err)
	}
	return sessions, nil
}

func (r *sessionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&domain.Session{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}
	return nil
}

func (r *sessionRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&domain.Session{}, "user_id = ?", userID).Error; err != nil {
		return fmt.Errorf("failed to delete user sessions: %w", err)
	}
	return nil
}

func (r *sessionRepository) DeleteExpired(ctx context.Context) error {
	if err := r.db.WithContext(ctx).Delete(&domain.Session{}, "expires_at < ?", time.Now()).Error; err != nil {
		return fmt.Errorf("failed to delete expired sessions: %w", err)
	}
	return nil
}

// OTP Repository
type OTPRepository interface {
	Create(ctx context.Context, otp *domain.OTP) error
	GetLatestByPhone(ctx context.Context, phone string) (*domain.OTP, error)
	MarkAsVerified(ctx context.Context, id uuid.UUID) error
	IncrementAttempts(ctx context.Context, id uuid.UUID) error
	DeleteExpired(ctx context.Context) error
}

type otpRepository struct {
	db *gorm.DB
}

func NewOTPRepository(db *gorm.DB) OTPRepository {
	return &otpRepository{db: db}
}

func (r *otpRepository) Create(ctx context.Context, otp *domain.OTP) error {
	if err := r.db.WithContext(ctx).Create(otp).Error; err != nil {
		return fmt.Errorf("failed to create OTP: %w", err)
	}
	return nil
}

func (r *otpRepository) GetLatestByPhone(ctx context.Context, phone string) (*domain.OTP, error) {
	var otp domain.OTP
	if err := r.db.WithContext(ctx).
		Where("phone = ? AND verified = false AND expires_at > ?", phone, time.Now()).
		Order("created_at DESC").
		First(&otp).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get OTP: %w", err)
	}
	return &otp, nil
}

func (r *otpRepository) MarkAsVerified(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Model(&domain.OTP{}).Where("id = ?", id).Update("verified", true).Error; err != nil {
		return fmt.Errorf("failed to mark OTP as verified: %w", err)
	}
	return nil
}

func (r *otpRepository) IncrementAttempts(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Model(&domain.OTP{}).Where("id = ?", id).UpdateColumn("attempts", gorm.Expr("attempts + ?", 1)).Error; err != nil {
		return fmt.Errorf("failed to increment OTP attempts: %w", err)
	}
	return nil
}

func (r *otpRepository) DeleteExpired(ctx context.Context) error {
	if err := r.db.WithContext(ctx).Delete(&domain.OTP{}, "expires_at < ?", time.Now()).Error; err != nil {
		return fmt.Errorf("failed to delete expired OTPs: %w", err)
	}
	return nil
}

// Password Reset Repository
type PasswordResetRepository interface {
	Create(ctx context.Context, reset *domain.PasswordReset) error
	GetByToken(ctx context.Context, token string) (*domain.PasswordReset, error)
	MarkAsUsed(ctx context.Context, id uuid.UUID) error
	DeleteExpired(ctx context.Context) error
}

type passwordResetRepository struct {
	db *gorm.DB
}

func NewPasswordResetRepository(db *gorm.DB) PasswordResetRepository {
	return &passwordResetRepository{db: db}
}

func (r *passwordResetRepository) Create(ctx context.Context, reset *domain.PasswordReset) error {
	if err := r.db.WithContext(ctx).Create(reset).Error; err != nil {
		return fmt.Errorf("failed to create password reset: %w", err)
	}
	return nil
}

func (r *passwordResetRepository) GetByToken(ctx context.Context, token string) (*domain.PasswordReset, error) {
	var reset domain.PasswordReset
	if err := r.db.WithContext(ctx).Preload("User").
		Where("token = ? AND used = false AND expires_at > ?", token, time.Now()).
		First(&reset).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get password reset: %w", err)
	}
	return &reset, nil
}

func (r *passwordResetRepository) MarkAsUsed(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Model(&domain.PasswordReset{}).Where("id = ?", id).Update("used", true).Error; err != nil {
		return fmt.Errorf("failed to mark password reset as used: %w", err)
	}
	return nil
}

func (r *passwordResetRepository) DeleteExpired(ctx context.Context) error {
	if err := r.db.WithContext(ctx).Delete(&domain.PasswordReset{}, "expires_at < ?", time.Now()).Error; err != nil {
		return fmt.Errorf("failed to delete expired password resets: %w", err)
	}
	return nil
}
