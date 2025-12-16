package domain

import (
	"time"

	"github.com/google/uuid"
)

type AddressType string

const (
	AddressTypeHome  AddressType = "home"
	AddressTypeWork  AddressType = "work"
	AddressTypeOther AddressType = "other"
)

type Address struct {
	ID            uuid.UUID   `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID        uuid.UUID   `gorm:"type:uuid;not null;index" json:"user_id"`
	Type          AddressType `gorm:"type:address_type;not null;default:'home'" json:"type"`
	FullName      string      `gorm:"size:255;not null" json:"full_name"`
	Phone         string      `gorm:"size:20;not null" json:"phone"`
	Email         *string     `gorm:"size:255" json:"email,omitempty"`
	StreetAddress string      `gorm:"size:255;not null" json:"street_address"`
	Landmark      *string     `gorm:"size:255" json:"landmark,omitempty"`
	City          string      `gorm:"size:100;not null" json:"city"`
	State         string      `gorm:"size:100;not null" json:"state"`
	Pincode       string      `gorm:"size:10;not null" json:"pincode"`
	Country       string      `gorm:"size:100;default:'India'" json:"country"`
	Latitude      *float64    `gorm:"type:decimal(10,8)" json:"latitude,omitempty"`
	Longitude     *float64    `gorm:"type:decimal(11,8)" json:"longitude,omitempty"`
	IsDefault     bool        `gorm:"default:false" json:"is_default"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
}

func (Address) TableName() string {
	return "addresses"
}
