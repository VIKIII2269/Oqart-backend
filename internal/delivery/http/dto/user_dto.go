package dto

import "time"

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
