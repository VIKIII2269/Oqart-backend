package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type ProductStatus string
type StockStatus string

const (
	ProductStatusDraft      ProductStatus = "draft"
	ProductStatusPending    ProductStatus = "pending"
	ProductStatusApproved   ProductStatus = "approved"
	ProductStatusRejected   ProductStatus = "rejected"
	ProductStatusActive     ProductStatus = "active"
	ProductStatusInactive   ProductStatus = "inactive"
	ProductStatusOutOfStock ProductStatus = "out_of_stock"

	StockStatusInStock      StockStatus = "in_stock"
	StockStatusLowStock     StockStatus = "low_stock"
	StockStatusOutOfStock   StockStatus = "out_of_stock"
	StockStatusDiscontinued StockStatus = "discontinued"
)

type Category struct {
	ID              uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	ParentID        *uuid.UUID `gorm:"type:uuid;index" json:"parent_id,omitempty"`
	Name            string     `gorm:"size:255;not null" json:"name"`
	Slug            string     `gorm:"size:255;not null;uniqueIndex" json:"slug"`
	Description     *string    `gorm:"type:text" json:"description,omitempty"`
	ImageURL        *string    `gorm:"size:500" json:"image_url,omitempty"`
	IconURL         *string    `gorm:"size:500" json:"icon_url,omitempty"`
	DisplayOrder    int        `gorm:"default:0" json:"display_order"`
	IsActive        bool       `gorm:"default:true" json:"is_active"`
	IsFeatured      bool       `gorm:"default:false" json:"is_featured"`
	MetaTitle       *string    `gorm:"size:255" json:"meta_title,omitempty"`
	MetaDescription *string    `gorm:"type:text" json:"meta_description,omitempty"`
	MetaKeywords    *string    `gorm:"type:text" json:"meta_keywords,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`

	Parent     *Category  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children   []Category `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}

type Product struct {
	ID                   uuid.UUID     `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	VendorID             uuid.UUID     `gorm:"type:uuid;not null;index" json:"vendor_id"`
	CategoryID           uuid.UUID     `gorm:"type:uuid;not null;index" json:"category_id"`
	SKU                  string        `gorm:"size:100;not null;uniqueIndex" json:"sku"`
	Name                 string        `gorm:"size:255;not null" json:"name"`
	Slug                 string        `gorm:"size:255;not null;uniqueIndex" json:"slug"`
	Description          *string       `gorm:"type:text" json:"description,omitempty"`
	ShortDescription     *string       `gorm:"type:text" json:"short_description,omitempty"`
	Status               ProductStatus `gorm:"type:product_status;default:'draft'" json:"status"`
	StockStatus          StockStatus   `gorm:"type:stock_status;default:'in_stock'" json:"stock_status"`
	Price                float64       `gorm:"type:decimal(15,2);not null" json:"price"`
	MRP                  float64       `gorm:"type:decimal(15,2);not null" json:"mrp"`
	DiscountPercentage   float64       `gorm:"type:decimal(5,2);default:0" json:"discount_percentage"`
	CostPrice            *float64      `gorm:"type:decimal(15,2)" json:"cost_price,omitempty"`
	StockQuantity        int           `gorm:"default:0" json:"stock_quantity"`
	LowStockThreshold    int           `gorm:"default:10" json:"low_stock_threshold"`
	AllowBackorder       bool          `gorm:"default:false" json:"allow_backorder"`
	Weight               *float64      `gorm:"type:decimal(10,2)" json:"weight,omitempty"`
	WeightUnit           string        `gorm:"size:20;default:'kg'" json:"weight_unit"`
	DimensionsLength     *float64      `gorm:"type:decimal(10,2)" json:"dimensions_length,omitempty"`
	DimensionsWidth      *float64      `gorm:"type:decimal(10,2)" json:"dimensions_width,omitempty"`
	DimensionsHeight     *float64      `gorm:"type:decimal(10,2)" json:"dimensions_height,omitempty"`
	DimensionsUnit       string        `gorm:"size:20;default:'cm'" json:"dimensions_unit"`
	Brand                *string       `gorm:"size:255" json:"brand,omitempty"`
	Manufacturer         *string       `gorm:"size:255" json:"manufacturer,omitempty"`
	CountryOfOrigin      string        `gorm:"size:100;default:'India'" json:"country_of_origin"`
	ShelfLife            *string       `gorm:"size:100" json:"shelf_life,omitempty"`
	Ingredients          *string       `gorm:"type:text" json:"ingredients,omitempty"`
	NutritionalInfo      pq.StringArray `gorm:"type:jsonb" json:"nutritional_info,omitempty"`
	UsageInstructions    *string       `gorm:"type:text" json:"usage_instructions,omitempty"`
	StorageInstructions  *string       `gorm:"type:text" json:"storage_instructions,omitempty"`
	IsFeatured           bool          `gorm:"default:false" json:"is_featured"`
	IsBestseller         bool          `gorm:"default:false" json:"is_bestseller"`
	IsNewArrival         bool          `gorm:"default:false" json:"is_new_arrival"`
	IsOrganic            bool          `gorm:"default:true" json:"is_organic"`
	IsVegan              bool          `gorm:"default:false" json:"is_vegan"`
	IsGlutenFree         bool          `gorm:"default:false" json:"is_gluten_free"`
	MetaTitle            *string       `gorm:"size:255" json:"meta_title,omitempty"`
	MetaDescription      *string       `gorm:"type:text" json:"meta_description,omitempty"`
	MetaKeywords         *string       `gorm:"type:text" json:"meta_keywords,omitempty"`
	ViewsCount           int           `gorm:"default:0" json:"views_count"`
	SalesCount           int           `gorm:"default:0" json:"sales_count"`
	Rating               float64       `gorm:"type:decimal(3,2);default:0" json:"rating"`
	ReviewCount          int           `gorm:"default:0" json:"review_count"`
	ApprovedAt           *time.Time    `json:"approved_at,omitempty"`
	ApprovedBy           *uuid.UUID    `gorm:"type:uuid" json:"approved_by,omitempty"`
	RejectionReason      *string       `gorm:"type:text" json:"rejection_reason,omitempty"`
	CreatedAt            time.Time     `json:"created_at"`
	UpdatedAt            time.Time     `json:"updated_at"`

	Category       Category             `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Images         []ProductImage       `gorm:"foreignKey:ProductID" json:"images,omitempty"`
	Variants       []ProductVariant     `gorm:"foreignKey:ProductID" json:"variants,omitempty"`
	Certifications []Certification      `gorm:"many2many:product_certifications" json:"certifications,omitempty"`
	Tags           []ProductTag         `gorm:"many2many:product_tag_mapping" json:"tags,omitempty"`
}

type ProductImage struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	ProductID    uuid.UUID `gorm:"type:uuid;not null;index" json:"product_id"`
	URL          string    `gorm:"size:500;not null" json:"url"`
	AltText      *string   `gorm:"size:255" json:"alt_text,omitempty"`
	DisplayOrder int       `gorm:"default:0" json:"display_order"`
	IsPrimary    bool      `gorm:"default:false" json:"is_primary"`
	CreatedAt    time.Time `json:"created_at"`
}

type ProductVariant struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	ProductID     uuid.UUID `gorm:"type:uuid;not null;index" json:"product_id"`
	SKU           string    `gorm:"size:100;not null;uniqueIndex" json:"sku"`
	Name          string    `gorm:"size:255;not null" json:"name"`
	Price         float64   `gorm:"type:decimal(15,2);not null" json:"price"`
	MRP           float64   `gorm:"type:decimal(15,2);not null" json:"mrp"`
	StockQuantity int       `gorm:"default:0" json:"stock_quantity"`
	Weight        *float64  `gorm:"type:decimal(10,2)" json:"weight,omitempty"`
	WeightUnit    string    `gorm:"size:20;default:'kg'" json:"weight_unit"`
	Attributes    pq.StringArray `gorm:"type:jsonb" json:"attributes,omitempty"`
	IsActive      bool      `gorm:"default:true" json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Certification struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name        string    `gorm:"size:255;not null;uniqueIndex" json:"name"`
	Slug        string    `gorm:"size:255;not null;uniqueIndex" json:"slug"`
	Description *string   `gorm:"type:text" json:"description,omitempty"`
	IconURL     *string   `gorm:"size:500" json:"icon_url,omitempty"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProductTag struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name      string    `gorm:"size:100;not null;uniqueIndex" json:"name"`
	Slug      string    `gorm:"size:100;not null;uniqueIndex" json:"slug"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName overrides
func (Category) TableName() string        { return "categories" }
func (Product) TableName() string         { return "products" }
func (ProductImage) TableName() string    { return "product_images" }
func (ProductVariant) TableName() string  { return "product_variants" }
func (Certification) TableName() string   { return "certifications" }
func (ProductTag) TableName() string      { return "product_tags" }
