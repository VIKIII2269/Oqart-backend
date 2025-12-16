package dto

import "time"

// UserResponse represents the user profile response
type UserResponse struct {
	ID            string     `json:"id"`
	Email         string     `json:"email"`
	Phone         *string    `json:"phone,omitempty"`
	FirstName     string     `json:"first_name"`
	LastName      string     `json:"last_name"`
	Role          string     `json:"role"`
	Status        string     `json:"status"`
	EmailVerified bool       `json:"email_verified"`
	PhoneVerified bool       `json:"phone_verified"`
	Avatar        *string    `json:"avatar,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	LastLoginAt   *time.Time `json:"last_login_at,omitempty"`
}

// UpdateProfileRequest represents profile update request
type UpdateProfileRequest struct {
	FirstName string  `json:"first_name" validate:"required,min=2,max=50"`
	LastName  string  `json:"last_name" validate:"required,min=2,max=50"`
	Phone     *string `json:"phone,omitempty" validate:"omitempty,indian_phone"`
	Avatar    *string `json:"avatar,omitempty" validate:"omitempty,url"`
}

// ChangePasswordRequest represents password change request
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,strong_password"`
}

// UpdateEmailRequest represents email update request
type UpdateEmailRequest struct {
	NewEmail string `json:"new_email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UpdatePhoneRequest represents phone update request
type UpdatePhoneRequest struct {
	NewPhone string `json:"new_phone" validate:"required,indian_phone"`
	Password string `json:"password" validate:"required"`
}

// PreferencesResponse represents user preferences
type PreferencesResponse struct {
	UserID                string    `json:"user_id"`
	Language              string    `json:"language"`
	Currency              string    `json:"currency"`
	NotificationEmail     bool      `json:"notification_email"`
	NotificationSMS       bool      `json:"notification_sms"`
	NotificationPush      bool      `json:"notification_push"`
	MarketingEmails       bool      `json:"marketing_emails"`
	OrderUpdates          bool      `json:"order_updates"`
	NewsletterSubscribed  bool      `json:"newsletter_subscribed"`
	TwoFactorEnabled      bool      `json:"two_factor_enabled"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

// UpdatePreferencesRequest represents preferences update request
type UpdatePreferencesRequest struct {
	Language             *string `json:"language,omitempty" validate:"omitempty,oneof=en hi"`
	Currency             *string `json:"currency,omitempty" validate:"omitempty,oneof=INR USD"`
	NotificationEmail    *bool   `json:"notification_email,omitempty"`
	NotificationSMS      *bool   `json:"notification_sms,omitempty"`
	NotificationPush     *bool   `json:"notification_push,omitempty"`
	MarketingEmails      *bool   `json:"marketing_emails,omitempty"`
	OrderUpdates         *bool   `json:"order_updates,omitempty"`
	NewsletterSubscribed *bool   `json:"newsletter_subscribed,omitempty"`
	TwoFactorEnabled     *bool   `json:"two_factor_enabled,omitempty"`
}

// AddressResponse represents address information
type AddressResponse struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	FullName     string    `json:"full_name"`
	Phone        string    `json:"phone"`
	AddressLine1 string    `json:"address_line1"`
	AddressLine2 *string   `json:"address_line2,omitempty"`
	City         string    `json:"city"`
	State        string    `json:"state"`
	Pincode      string    `json:"pincode"`
	Country      string    `json:"country"`
	AddressType  string    `json:"address_type"`
	IsDefault    bool      `json:"is_default"`
	Latitude     *float64  `json:"latitude,omitempty"`
	Longitude    *float64  `json:"longitude,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CreateAddressRequest represents address creation request
type CreateAddressRequest struct {
	FullName     string   `json:"full_name" validate:"required,min=3,max=100"`
	Phone        string   `json:"phone" validate:"required,indian_phone"`
	AddressLine1 string   `json:"address_line1" validate:"required,min=5,max=255"`
	AddressLine2 *string  `json:"address_line2,omitempty" validate:"omitempty,max=255"`
	City         string   `json:"city" validate:"required,min=2,max=100"`
	State        string   `json:"state" validate:"required,min=2,max=100"`
	Pincode      string   `json:"pincode" validate:"required,indian_pincode"`
	Country      string   `json:"country" validate:"required,min=2,max=100"`
	AddressType  string   `json:"address_type" validate:"required,oneof=home work other"`
	IsDefault    bool     `json:"is_default"`
	Latitude     *float64 `json:"latitude,omitempty" validate:"omitempty,latitude"`
	Longitude    *float64 `json:"longitude,omitempty" validate:"omitempty,longitude"`
}

// UpdateAddressRequest represents address update request
type UpdateAddressRequest struct {
	FullName     *string  `json:"full_name,omitempty" validate:"omitempty,min=3,max=100"`
	Phone        *string  `json:"phone,omitempty" validate:"omitempty,indian_phone"`
	AddressLine1 *string  `json:"address_line1,omitempty" validate:"omitempty,min=5,max=255"`
	AddressLine2 *string  `json:"address_line2,omitempty" validate:"omitempty,max=255"`
	City         *string  `json:"city,omitempty" validate:"omitempty,min=2,max=100"`
	State        *string  `json:"state,omitempty" validate:"omitempty,min=2,max=100"`
	Pincode      *string  `json:"pincode,omitempty" validate:"omitempty,indian_pincode"`
	Country      *string  `json:"country,omitempty" validate:"omitempty,min=2,max=100"`
	AddressType  *string  `json:"address_type,omitempty" validate:"omitempty,oneof=home work other"`
	Latitude     *float64 `json:"latitude,omitempty" validate:"omitempty,latitude"`
	Longitude    *float64 `json:"longitude,omitempty" validate:"omitempty,longitude"`
}
