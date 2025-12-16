package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string
type UserStatus string
type AuthProvider string

const (
	RoleCustomer   UserRole = "customer"
	RoleVendor     UserRole = "vendor"
	RoleAdmin      UserRole = "admin"
	RoleSuperAdmin UserRole = "super_admin"

	StatusActive    UserStatus = "active"
	StatusInactive  UserStatus = "inactive"
	StatusSuspended UserStatus = "suspended"
	StatusDeleted   UserStatus = "deleted"

	AuthProviderEmail    AuthProvider = "email"
	AuthProviderPhone    AuthProvider = "phone"
	AuthProviderGoogle   AuthProvider = "google"
	AuthProviderFacebook AuthProvider = "facebook"
)

type User struct {
	ID             uuid.UUID    `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Email          *string      `gorm:"uniqueIndex;size:255" json:"email,omitempty"`
	Phone          *string      `gorm:"uniqueIndex;size:20" json:"phone,omitempty"`
	PasswordHash   *string      `gorm:"size:255" json:"-"`
	Role           UserRole     `gorm:"type:user_role;default:'customer'" json:"role"`
	Status         UserStatus   `gorm:"type:user_status;default:'active'" json:"status"`
	AuthProvider   AuthProvider `gorm:"type:auth_provider;default:'email'" json:"auth_provider"`
	EmailVerified  bool         `gorm:"default:false" json:"email_verified"`
	PhoneVerified  bool         `gorm:"default:false" json:"phone_verified"`
	FirstName      *string      `gorm:"size:100" json:"first_name,omitempty"`
	LastName       *string      `gorm:"size:100" json:"last_name,omitempty"`
	AvatarURL      *string      `gorm:"size:500" json:"avatar_url,omitempty"`
	DateOfBirth    *time.Time   `json:"date_of_birth,omitempty"`
	Gender         *string      `gorm:"size:20" json:"gender,omitempty"`
	LastLoginAt    *time.Time   `json:"last_login_at,omitempty"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}

type Session struct {
	ID                uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID            uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	AccessToken       string    `gorm:"type:text;not null" json:"access_token"`
	RefreshToken      string    `gorm:"type:text;not null" json:"refresh_token"`
	UserAgent         *string   `gorm:"size:500" json:"user_agent,omitempty"`
	IPAddress         *string   `gorm:"size:45" json:"ip_address,omitempty"`
	ExpiresAt         time.Time `json:"expires_at"`
	RefreshExpiresAt  time.Time `json:"refresh_expires_at"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
}

type PasswordReset struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Token     string    `gorm:"size:255;not null;uniqueIndex" json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	Used      bool      `gorm:"default:false" json:"used"`
	CreatedAt time.Time `json:"created_at"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
}

type EmailVerification struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Token     string    `gorm:"size:255;not null;uniqueIndex" json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	Verified  bool      `gorm:"default:false" json:"verified"`
	CreatedAt time.Time `json:"created_at"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
}

type OTP struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Phone     string    `gorm:"size:20;not null;index" json:"phone"`
	Code      string    `gorm:"size:10;not null" json:"code"`
	ExpiresAt time.Time `json:"expires_at"`
	Verified  bool      `gorm:"default:false" json:"verified"`
	Attempts  int       `gorm:"default:0" json:"attempts"`
	CreatedAt time.Time `json:"created_at"`
}

type UserPreferences struct {
	ID                   uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID               uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`
	Language             string    `gorm:"size:10;default:'en'" json:"language"`
	Currency             string    `gorm:"size:10;default:'INR'" json:"currency"`
	Timezone             string    `gorm:"size:50;default:'Asia/Kolkata'" json:"timezone"`
	NotificationEmail    bool      `gorm:"default:true" json:"notification_email"`
	NotificationSMS      bool      `gorm:"default:true" json:"notification_sms"`
	NotificationPush     bool      `gorm:"default:true" json:"notification_push"`
	MarketingEmails      bool      `gorm:"default:false" json:"marketing_emails"`
	OrderUpdates         bool      `gorm:"default:true" json:"order_updates"`
	NewsletterSubscribed bool      `gorm:"default:false" json:"newsletter_subscribed"`
	TwoFactorEnabled     bool      `gorm:"default:false" json:"two_factor_enabled"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
}

type Address struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID     uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Type       string    `gorm:"size:50;default:'home'" json:"type"`
	FullName   string    `gorm:"size:255;not null" json:"full_name"`
	Phone      string    `gorm:"size:20;not null" json:"phone"`
	Street     string    `gorm:"size:255;not null" json:"street"`
	Landmark   *string   `gorm:"size:255" json:"landmark,omitempty"`
	City       string    `gorm:"size:100;not null" json:"city"`
	State      string    `gorm:"size:100;not null" json:"state"`
	Pincode    string    `gorm:"size:10;not null" json:"pincode"`
	Country    string    `gorm:"size:100;default:'India'" json:"country"`
	IsDefault  bool      `gorm:"default:false" json:"is_default"`
	Latitude   *float64  `gorm:"type:decimal(10,8)" json:"latitude,omitempty"`
	Longitude  *float64  `gorm:"type:decimal(11,8)" json:"longitude,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
}

// TableName overrides
func (User) TableName() string                 { return "users" }
func (Session) TableName() string              { return "sessions" }
func (PasswordReset) TableName() string        { return "password_resets" }
func (EmailVerification) TableName() string    { return "email_verifications" }
func (OTP) TableName() string                  { return "otps" }
func (UserPreferences) TableName() string      { return "user_preferences" }
func (Address) TableName() string              { return "addresses" }

// Helper methods
func (u *User) FullName() string {
	if u.FirstName != nil && u.LastName != nil {
		return *u.FirstName + " " + *u.LastName
	}
	if u.FirstName != nil {
		return *u.FirstName
	}
	if u.LastName != nil {
		return *u.LastName
	}
	return ""
}

func (u *User) IsActive() bool {
	return u.Status == StatusActive
}

func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin || u.Role == RoleSuperAdmin
}

func (u *User) IsVendor() bool {
	return u.Role == RoleVendor
}

func (u *User) IsCustomer() bool {
	return u.Role == RoleCustomer
}
