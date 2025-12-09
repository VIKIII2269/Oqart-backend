package validator

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate
)

func init() {
	validate = validator.New()

	// Register custom validators
	_ = validate.RegisterValidation("phone", validatePhone)
	_ = validate.RegisterValidation("gstin", validateGSTIN)
	_ = validate.RegisterValidation("pan", validatePAN)
	_ = validate.RegisterValidation("ifsc", validateIFSC)
	_ = validate.RegisterValidation("pincode", validatePincode)
	_ = validate.RegisterValidation("strong_password", validateStrongPassword)
}

// Validate validates a struct
func Validate(data interface{}) error {
	if err := validate.Struct(data); err != nil {
		return formatValidationError(err)
	}
	return nil
}

// ValidateVar validates a single variable
func ValidateVar(field interface{}, tag string) error {
	if err := validate.Var(field, tag); err != nil {
		return formatValidationError(err)
	}
	return nil
}

// formatValidationError formats validation errors into a readable message
func formatValidationError(err error) error {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var messages []string
		for _, e := range validationErrors {
			messages = append(messages, formatFieldError(e))
		}
		return fmt.Errorf(strings.Join(messages, "; "))
	}
	return err
}

// formatFieldError formats a single field error
func formatFieldError(e validator.FieldError) string {
	field := strings.ToLower(e.Field())

	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", field, e.Param())
	case "max":
		return fmt.Sprintf("%s must not exceed %s characters", field, e.Param())
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters", field, e.Param())
	case "phone":
		return fmt.Sprintf("%s must be a valid Indian phone number", field)
	case "gstin":
		return fmt.Sprintf("%s must be a valid GSTIN", field)
	case "pan":
		return fmt.Sprintf("%s must be a valid PAN", field)
	case "ifsc":
		return fmt.Sprintf("%s must be a valid IFSC code", field)
	case "pincode":
		return fmt.Sprintf("%s must be a valid Indian pincode", field)
	case "strong_password":
		return fmt.Sprintf("%s must be at least 8 characters and contain uppercase, lowercase, number, and special character", field)
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID", field)
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	case "numeric":
		return fmt.Sprintf("%s must be numeric", field)
	case "alpha":
		return fmt.Sprintf("%s must contain only alphabetic characters", field)
	case "alphanum":
		return fmt.Sprintf("%s must contain only alphanumeric characters", field)
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, e.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", field, e.Param())
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", field, e.Param())
	case "lt":
		return fmt.Sprintf("%s must be less than %s", field, e.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, e.Param())
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}

// Custom validators

// validatePhone validates Indian phone numbers
func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	// Indian phone number: 10 digits, optionally starting with +91 or 91
	pattern := `^(\+91|91)?[6-9]\d{9}$`
	matched, _ := regexp.MatchString(pattern, phone)
	return matched
}

// validateGSTIN validates GSTIN format
func validateGSTIN(fl validator.FieldLevel) bool {
	gstin := fl.Field().String()
	// GSTIN format: 15 characters
	// Format: 22AAAAA0000A1Z5
	pattern := `^[0-9]{2}[A-Z]{5}[0-9]{4}[A-Z]{1}[1-9A-Z]{1}Z[0-9A-Z]{1}$`
	matched, _ := regexp.MatchString(pattern, gstin)
	return matched
}

// validatePAN validates PAN format
func validatePAN(fl validator.FieldLevel) bool {
	pan := fl.Field().String()
	// PAN format: 10 characters
	// Format: AAAAA9999A
	pattern := `^[A-Z]{5}[0-9]{4}[A-Z]{1}$`
	matched, _ := regexp.MatchString(pattern, pan)
	return matched
}

// validateIFSC validates IFSC code format
func validateIFSC(fl validator.FieldLevel) bool {
	ifsc := fl.Field().String()
	// IFSC format: 11 characters
	// Format: AAAA0BBBBBB
	pattern := `^[A-Z]{4}0[A-Z0-9]{6}$`
	matched, _ := regexp.MatchString(pattern, ifsc)
	return matched
}

// validatePincode validates Indian pincode
func validatePincode(fl validator.FieldLevel) bool {
	pincode := fl.Field().String()
	// Indian pincode: 6 digits
	pattern := `^[1-9][0-9]{5}$`
	matched, _ := regexp.MatchString(pattern, pincode)
	return matched
}

// validateStrongPassword validates password strength
func validateStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 8 {
		return false
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case 'a' <= char && char <= 'z':
			hasLower = true
		case '0' <= char && char <= '9':
			hasNumber = true
		case strings.ContainsRune("!@#$%^&*()_+-=[]{}|;:,.<>?", char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}

// Helper functions

// ValidateEmail validates email format
func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

// ValidatePhone validates phone number
func ValidatePhone(phone string) bool {
	pattern := `^(\+91|91)?[6-9]\d{9}$`
	matched, _ := regexp.MatchString(pattern, phone)
	return matched
}

// ValidateGSTIN validates GSTIN
func ValidateGSTIN(gstin string) bool {
	pattern := `^[0-9]{2}[A-Z]{5}[0-9]{4}[A-Z]{1}[1-9A-Z]{1}Z[0-9A-Z]{1}$`
	matched, _ := regexp.MatchString(pattern, gstin)
	return matched
}

// ValidatePAN validates PAN
func ValidatePAN(pan string) bool {
	pattern := `^[A-Z]{5}[0-9]{4}[A-Z]{1}$`
	matched, _ := regexp.MatchString(pattern, pan)
	return matched
}

// ValidateIFSC validates IFSC code
func ValidateIFSC(ifsc string) bool {
	pattern := `^[A-Z]{4}0[A-Z0-9]{6}$`
	matched, _ := regexp.MatchString(pattern, ifsc)
	return matched
}

// ValidatePincode validates pincode
func ValidatePincode(pincode string) bool {
	pattern := `^[1-9][0-9]{5}$`
	matched, _ := regexp.MatchString(pattern, pincode)
	return matched
}
