package dto

import (
	"time"

	"github.com/oqart/backend/internal/domain"
)

// Register Request
type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Phone     string `json:"phone" validate:"omitempty,phone"`
	Password  string `json:"password" validate:"required,strong_password"`
	FirstName string `json:"first_name" validate:"required,min=2,max=100"`
	LastName  string `json:"last_name" validate:"omitempty,max=100"`
}

// Login Request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// Phone Login Request (Send OTP)
type PhoneLoginRequest struct {
	Phone string `json:"phone" validate:"required,phone"`
}

// Verify OTP Request
type VerifyOTPRequest struct {
	Phone string `json:"phone" validate:"required,phone"`
	Code  string `json:"code" validate:"required,len=6,numeric"`
}

// Google OAuth Request
type GoogleAuthRequest struct {
	IDToken string `json:"id_token" validate:"required"`
}

// Refresh Token Request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// Forgot Password Request
type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// Reset Password Request
type ResetPasswordRequest struct {
	Token    string `json:"token" validate:"required"`
	Password string `json:"password" validate:"required,strong_password"`
}

// Verify Email Request
type VerifyEmailRequest struct {
	Token string `json:"token" validate:"required"`
}

// Change Password Request
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,strong_password"`
}

// Auth Response
type AuthResponse struct {
	User         *UserResponse `json:"user"`
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	ExpiresAt    time.Time     `json:"expires_at"`
}

// User Response
type UserResponse struct {
	ID            string     `json:"id"`
	Email         *string    `json:"email,omitempty"`
	Phone         *string    `json:"phone,omitempty"`
	Role          string     `json:"role"`
	Status        string     `json:"status"`
	EmailVerified bool       `json:"email_verified"`
	PhoneVerified bool       `json:"phone_verified"`
	FirstName     *string    `json:"first_name,omitempty"`
	LastName      *string    `json:"last_name,omitempty"`
	AvatarURL     *string    `json:"avatar_url,omitempty"`
	DateOfBirth   *time.Time `json:"date_of_birth,omitempty"`
	Gender        *string    `json:"gender,omitempty"`
	LastLoginAt   *time.Time `json:"last_login_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// Convert domain user to DTO
func ToUserResponse(user *domain.User) *UserResponse {
	return &UserResponse{
		ID:            user.ID.String(),
		Email:         user.Email,
		Phone:         user.Phone,
		Role:          string(user.Role),
		Status:        string(user.Status),
		EmailVerified: user.EmailVerified,
		PhoneVerified: user.PhoneVerified,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		AvatarURL:     user.AvatarURL,
		DateOfBirth:   user.DateOfBirth,
		Gender:        user.Gender,
		LastLoginAt:   user.LastLoginAt,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
	}
}

// Update Profile Request
type UpdateProfileRequest struct {
	FirstName   *string    `json:"first_name" validate:"omitempty,min=2,max=100"`
	LastName    *string    `json:"last_name" validate:"omitempty,max=100"`
	DateOfBirth *time.Time `json:"date_of_birth" validate:"omitempty"`
	Gender      *string    `json:"gender" validate:"omitempty,oneof=male female other"`
}

// Update Email Request
type UpdateEmailRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// Update Phone Request
type UpdatePhoneRequest struct {
	Phone string `json:"phone" validate:"required,phone"`
}

// OTP Response
type OTPResponse struct {
	Message   string    `json:"message"`
	ExpiresAt time.Time `json:"expires_at"`
}

// Message Response
type MessageResponse struct {
	Message string `json:"message"`
}
