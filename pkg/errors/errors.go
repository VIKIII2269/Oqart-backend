package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// AppError represents an application error
type AppError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
	Err        error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// New creates a new AppError
func New(code, message string, statusCode int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
	}
}

// Wrap wraps an error with additional context
func Wrap(err error, code, message string, statusCode int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
		Err:        err,
	}
}

// Common errors
var (
	// 400 Bad Request
	ErrBadRequest       = New("BAD_REQUEST", "Bad request", http.StatusBadRequest)
	ErrInvalidInput     = New("INVALID_INPUT", "Invalid input", http.StatusBadRequest)
	ErrValidation       = New("VALIDATION_ERROR", "Validation error", http.StatusBadRequest)
	ErrInvalidCredentials = New("INVALID_CREDENTIALS", "Invalid credentials", http.StatusBadRequest)
	ErrInvalidToken     = New("INVALID_TOKEN", "Invalid token", http.StatusBadRequest)

	// 401 Unauthorized
	ErrUnauthorized    = New("UNAUTHORIZED", "Unauthorized", http.StatusUnauthorized)
	ErrTokenExpired    = New("TOKEN_EXPIRED", "Token expired", http.StatusUnauthorized)
	ErrAuthRequired    = New("AUTH_REQUIRED", "Authentication required", http.StatusUnauthorized)

	// 403 Forbidden
	ErrForbidden        = New("FORBIDDEN", "Forbidden", http.StatusForbidden)
	ErrPermissionDenied = New("PERMISSION_DENIED", "Permission denied", http.StatusForbidden)

	// 404 Not Found
	ErrNotFound        = New("NOT_FOUND", "Resource not found", http.StatusNotFound)
	ErrUserNotFound    = New("USER_NOT_FOUND", "User not found", http.StatusNotFound)
	ErrProductNotFound = New("PRODUCT_NOT_FOUND", "Product not found", http.StatusNotFound)
	ErrOrderNotFound   = New("ORDER_NOT_FOUND", "Order not found", http.StatusNotFound)

	// 409 Conflict
	ErrConflict         = New("CONFLICT", "Resource conflict", http.StatusConflict)
	ErrEmailExists      = New("EMAIL_EXISTS", "Email already exists", http.StatusConflict)
	ErrPhoneExists      = New("PHONE_EXISTS", "Phone number already exists", http.StatusConflict)
	ErrDuplicateEntry   = New("DUPLICATE_ENTRY", "Duplicate entry", http.StatusConflict)

	// 422 Unprocessable Entity
	ErrUnprocessable    = New("UNPROCESSABLE", "Unprocessable entity", http.StatusUnprocessableEntity)
	ErrOutOfStock       = New("OUT_OF_STOCK", "Product out of stock", http.StatusUnprocessableEntity)
	ErrInsufficientStock = New("INSUFFICIENT_STOCK", "Insufficient stock", http.StatusUnprocessableEntity)

	// 500 Internal Server Error
	ErrInternal          = New("INTERNAL_ERROR", "Internal server error", http.StatusInternalServerError)
	ErrDatabase          = New("DATABASE_ERROR", "Database error", http.StatusInternalServerError)
	ErrCache             = New("CACHE_ERROR", "Cache error", http.StatusInternalServerError)

	// 503 Service Unavailable
	ErrServiceUnavailable = New("SERVICE_UNAVAILABLE", "Service unavailable", http.StatusServiceUnavailable)
)

// HTTP response helper
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

func NewErrorResponse(err error) *ErrorResponse {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return &ErrorResponse{
			Code:    appErr.Code,
			Message: appErr.Message,
		}
	}

	return &ErrorResponse{
		Code:    "INTERNAL_ERROR",
		Message: "Internal server error",
	}
}

func NewErrorResponseWithDetails(err error, details interface{}) *ErrorResponse {
	resp := NewErrorResponse(err)
	resp.Details = details
	return resp
}

func GetStatusCode(err error) int {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.StatusCode
	}
	return http.StatusInternalServerError
}
