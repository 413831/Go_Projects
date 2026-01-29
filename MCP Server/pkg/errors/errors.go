package errors

import (
	"fmt"
)

// AppError represents an application error
type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("[%s] %s: %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// NewAppError creates a new application error
func NewAppError(code, message, details string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// Predefined error codes
const (
	ErrCodeValidation   = "VALIDATION_ERROR"
	ErrCodeGitHubAPI    = "GITHUB_API_ERROR"
	ErrCodeNetwork      = "NETWORK_ERROR"
	ErrCodeJSONDecoding = "JSON_DECODING_ERROR"
	ErrCodeUnauthorized = "UNAUTHORIZED"
	ErrCodeNotFound     = "NOT_FOUND"
)

// NewValidationError creates a validation AppError
func NewValidationError(message string) *AppError {
	return NewAppError(ErrCodeValidation, message, "")
}

func NewGitHubAPIError(message string) *AppError {
	return NewAppError(ErrCodeGitHubAPI, message, "")
}

func NewNetworkError(details string) *AppError {
	return NewAppError(ErrCodeNetwork, "Network error", details)
}

func NewJSONDecodingError(details string) *AppError {
	return NewAppError(ErrCodeJSONDecoding, "Error decoding JSON", details)
}

func NewUnauthorizedError() *AppError {
	return NewAppError(ErrCodeUnauthorized, "Unauthorized", "GitHub token is invalid or missing")
}

func NewNotFoundError(resource string) *AppError {
	return NewAppError(ErrCodeNotFound, "Resource not found", resource)
}
