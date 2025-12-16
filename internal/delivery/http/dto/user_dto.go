package dto

import (
	"time"

	"github.com/oqart/backend/internal/domain"
)

// User DTOs are already defined in auth_dto.go (UserResponse)
// Adding additional DTOs for user management

// Address DTOs
type AddressRequest struct {
	Type          string  `json:"type" validate:"required,oneof=home work other"`
	FullName      string  `json:"full_name" validate:"required,min=2,max=255"`
	Phone         string  `json:"phone" validate:"required,phone"`
	Email         string  `json:"email" validate:"omitempty,email"`
	StreetAddress string  `json:"street_address" validate:"required,max=255"`
	Landmark      string  `json:"landmark" validate:"omitempty,max=255"`
	City          string  `json:"city" validate:"required,max=100"`
	State         string  `json:"state" validate:"required,max=100"`
	Pincode       string  `json:"pincode" validate:"required,pincode"`
	Country       string  `json:"country" validate:"omitempty,max=100"`
	IsDefault     bool    `json:"is_default"`
}

type AddressResponse struct {
	ID            string    `json:"id"`
	Type          string    `json:"type"`
	FullName      string    `json:"full_name"`
	Phone         string    `json:"phone"`
	Email         *string   `json:"email,omitempty"`
	StreetAddress string    `json:"street_address"`
	Landmark      *string   `json:"landmark,omitempty"`
	City          string    `json:"city"`
	State         string    `json:"state"`
	Pincode       string    `json:"pincode"`
	Country       string    `json:"country"`
	IsDefault     bool      `json:"is_default"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// User Preferences DTOs
type UserPreferencesResponse struct {
	Language           string `json:"language"`
	Currency           string `json:"currency"`
	Timezone           string `json:"timezone"`
	EmailNotifications bool   `json:"email_notifications"`
	SMSNotifications   bool   `json:"sms_notifications"`
	PushNotifications  bool   `json:"push_notifications"`
	MarketingEmails    bool   `json:"marketing_emails"`
}

type UpdatePreferencesRequest struct {
	Language           *string `json:"language" validate:"omitempty,min=2,max=10"`
	Currency           *string `json:"currency" validate:"omitempty,min=3,max=10"`
	Timezone           *string `json:"timezone" validate:"omitempty,max=50"`
	EmailNotifications *bool   `json:"email_notifications"`
	SMSNotifications   *bool   `json:"sms_notifications"`
	PushNotifications  *bool   `json:"push_notifications"`
	MarketingEmails    *bool   `json:"marketing_emails"`
}

// Convert domain address to DTO
func ToAddressResponse(address *domain.Address) *AddressResponse {
	return &AddressResponse{
		ID:            address.ID.String(),
		Type:          string(address.Type),
		FullName:      address.FullName,
		Phone:         address.Phone,
		Email:         address.Email,
		StreetAddress: address.StreetAddress,
		Landmark:      address.Landmark,
		City:          address.City,
		State:         address.State,
		Pincode:       address.Pincode,
		Country:       address.Country,
		IsDefault:     address.IsDefault,
		CreatedAt:     address.CreatedAt,
		UpdatedAt:     address.UpdatedAt,
	}
}

// Convert domain user preferences to DTO
func ToUserPreferencesResponse(prefs *domain.UserPreferences) *UserPreferencesResponse {
	return &UserPreferencesResponse{
		Language:           prefs.Language,
		Currency:           prefs.Currency,
		Timezone:           prefs.Timezone,
		EmailNotifications: prefs.EmailNotifications,
		SMSNotifications:   prefs.SMSNotifications,
		PushNotifications:  prefs.PushNotifications,
		MarketingEmails:    prefs.MarketingEmails,
	}
}
