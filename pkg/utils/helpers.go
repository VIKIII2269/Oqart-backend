package utils

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
)

// GenerateUUID generates a new UUID
func GenerateUUID() string {
	return uuid.New().String()
}

// ParseUUID parses a UUID string
func ParseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

// IsValidUUID checks if a string is a valid UUID
func IsValidUUID(s string) bool {
	_, err := uuid.Parse(s)
	return err == nil
}

// GenerateSlug generates a URL-friendly slug from a string
func GenerateSlug(s string) string {
	// Convert to lowercase
	s = strings.ToLower(s)

	// Remove accents and special characters
	s = removeAccents(s)

	// Replace spaces and underscores with hyphens
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "_", "-")

	// Remove all non-alphanumeric characters except hyphens
	reg := regexp.MustCompile("[^a-z0-9-]+")
	s = reg.ReplaceAllString(s, "")

	// Replace multiple hyphens with single hyphen
	reg = regexp.MustCompile("-+")
	s = reg.ReplaceAllString(s, "-")

	// Trim hyphens from start and end
	s = strings.Trim(s, "-")

	return s
}

// removeAccents removes accents from characters
func removeAccents(s string) string {
	var b strings.Builder
	for _, r := range s {
		if unicode.In(r, unicode.Latin) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// GenerateUniqueSlug generates a unique slug with timestamp
func GenerateUniqueSlug(s string) string {
	slug := GenerateSlug(s)
	timestamp := time.Now().Unix()
	return fmt.Sprintf("%s-%d", slug, timestamp)
}

// Truncate truncates a string to a specified length
func Truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// Contains checks if a slice contains a string
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// ContainsInt checks if a slice contains an int
func ContainsInt(slice []int, item int) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// Remove removes an item from a slice
func Remove(slice []string, item string) []string {
	var result []string
	for _, s := range slice {
		if s != item {
			result = append(result, s)
		}
	}
	return result
}

// Unique returns unique elements from a slice
func Unique(slice []string) []string {
	keys := make(map[string]bool)
	var result []string
	for _, entry := range slice {
		if _, exists := keys[entry]; !exists {
			keys[entry] = true
			result = append(result, entry)
		}
	}
	return result
}

// SanitizePhone removes common phone number formatting
func SanitizePhone(phone string) string {
	// Remove spaces, hyphens, parentheses
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	phone = strings.ReplaceAll(phone, "(", "")
	phone = strings.ReplaceAll(phone, ")", "")
	phone = strings.ReplaceAll(phone, "+", "")

	// If starts with 91, keep it
	if strings.HasPrefix(phone, "91") && len(phone) == 12 {
		return phone
	}

	// If starts with 0, remove it
	if strings.HasPrefix(phone, "0") {
		phone = phone[1:]
	}

	return phone
}

// FormatPhone formats a phone number for display
func FormatPhone(phone string) string {
	phone = SanitizePhone(phone)

	if len(phone) == 10 {
		return fmt.Sprintf("+91 %s %s", phone[:5], phone[5:])
	}

	return phone
}

// FormatCurrency formats a number as currency
func FormatCurrency(amount float64) string {
	return fmt.Sprintf("₹%.2f", amount)
}

// Pointer returns a pointer to the value
func Pointer[T any](value T) *T {
	return &value
}

// StringPointer returns a pointer to a string
func StringPointer(s string) *string {
	return &s
}

// IntPointer returns a pointer to an int
func IntPointer(i int) *int {
	return &i
}

// BoolPointer returns a pointer to a bool
func BoolPointer(b bool) *bool {
	return &b
}

// DefaultString returns default value if string is empty
func DefaultString(s, defaultVal string) string {
	if s == "" {
		return defaultVal
	}
	return s
}

// DefaultInt returns default value if int is zero
func DefaultInt(i, defaultVal int) int {
	if i == 0 {
		return defaultVal
	}
	return i
}

// Pagination helper
type PaginationParams struct {
	Page     int
	PageSize int
	Offset   int
}

func NewPaginationParams(page, pageSize int) PaginationParams {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	return PaginationParams{
		Page:     page,
		PageSize: pageSize,
		Offset:   offset,
	}
}

type PaginationResponse struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

func NewPaginationResponse(page, pageSize, totalItems int) PaginationResponse {
	totalPages := (totalItems + pageSize - 1) / pageSize
	return PaginationResponse{
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}
}
