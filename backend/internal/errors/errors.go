package errors

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Standard error codes
const (
	ErrorCodeValidation      = "VALIDATION_ERROR"
	ErrorCodeUnauthorized    = "UNAUTHORIZED"
	ErrorCodeForbidden       = "FORBIDDEN"
	ErrorCodeNotFound        = "NOT_FOUND"
	ErrorCodeUsernameTaken   = "USERNAME_TAKEN"
	ErrorCodeVersionConflict = "VERSION_CONFLICT"
	ErrorCodeInternal        = "INTERNAL_ERROR"
)

// AppError represents a structured application error.
type AppError struct {
	Code      int    `json:"code"`
	ErrorCode string `json:"error_code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
	Details   string `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(code int, errorCode, message string) *AppError {
	return &AppError{
		Code:      code,
		ErrorCode: errorCode,
		Message:   message,
	}
}

func NewValidationError(message string) *AppError {
	return NewAppError(http.StatusBadRequest, ErrorCodeValidation, message)
}

func NewUnauthorizedError(message string) *AppError {
	return NewAppError(http.StatusUnauthorized, ErrorCodeUnauthorized, message)
}

func NewNotFoundError(message string) *AppError {
	return NewAppError(http.StatusNotFound, ErrorCodeNotFound, message)
}

func NewConflictError(errorCode, message string) *AppError {
	return NewAppError(http.StatusConflict, errorCode, message)
}

func NewVersionConflictError(details string) *AppError {
	err := NewAppError(http.StatusConflict, ErrorCodeVersionConflict, "resource conflict due to version mismatch")
	if details != "" {
		err.Details = details
	}
	return err
}

func NewInternalError() *AppError {
	return NewAppError(http.StatusInternalServerError, ErrorCodeInternal, "internal server error")
}

// NewValidationErrorFromBinding extracts validation errors from ShouldBindJSON errors
// and returns a descriptive AppError with field-level details.
func NewValidationErrorFromBinding(bindErr error) *AppError {
	var ve validator.ValidationErrors
	if ok := AsValidationErrors(bindErr, &ve); ok && len(ve) > 0 {
		var errMsgs []string
		for _, fe := range ve {
			errMsgs = append(errMsgs, fmt.Sprintf("field '%s' %s", fe.Field(), fe.Tag()))
		}
		msg := "validation failed: " + strings.Join(errMsgs, "; ")
		return NewAppError(http.StatusBadRequest, ErrorCodeValidation, msg)
	}
	return NewValidationError("invalid request body")
}

// AsValidationErrors checks if the error is a validator.ValidationErrors.
// Exposed as a function to avoid import cycle issues in handlers.
var AsValidationErrors = func(err error, target *validator.ValidationErrors) bool {
	if err == nil {
		return false
	}
	ve, ok := err.(validator.ValidationErrors)
	if ok {
		*target = ve
		return true
	}
	// Also check wrapped errors
	if errors.As(err, &ve) {
		*target = ve
		return true
	}
	return false
}
