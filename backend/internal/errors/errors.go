package errors

import (
	"encoding/json"
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
	ErrorCodeTokenExpired    = "TOKEN_EXPIRED"
	ErrorCodeInvalidToken    = "INVALID_TOKEN"
	ErrorCodeForbidden       = "FORBIDDEN"
	ErrorCodeNotFound        = "NOT_FOUND"
	ErrorCodeUsernameTaken   = "USERNAME_TAKEN"
	ErrorCodeVersionConflict = "VERSION_CONFLICT"
	ErrorCodeInternal        = "INTERNAL_ERROR"
)

// AppError represents a structured application error.
type AppError struct {
	Code      int         `json:"code"`
	ErrorCode string      `json:"error_code"`
	Message   string      `json:"message"`
	RequestID string      `json:"request_id"`
	Details   interface{} `json:"details,omitempty"`
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

func NewTokenExpiredError() *AppError {
	return NewAppError(http.StatusUnauthorized, ErrorCodeTokenExpired, "token has expired")
}

func NewInvalidTokenError() *AppError {
	return NewAppError(http.StatusUnauthorized, ErrorCodeInvalidToken, "invalid token")
}

func NewNotFoundError(message string) *AppError {
	return NewAppError(http.StatusNotFound, ErrorCodeNotFound, message)
}

func NewConflictError(errorCode, message string) *AppError {
	return NewAppError(http.StatusConflict, errorCode, message)
}

func NewVersionConflictError(currentVersion int64) *AppError {
	err := NewAppError(http.StatusConflict, ErrorCodeVersionConflict, "resource conflict due to version mismatch")
	if currentVersion > 0 {
		err.Details = map[string]interface{}{
			"current_version": currentVersion,
		}
	}
	return err
}

func NewInternalError() *AppError {
	return NewAppError(http.StatusInternalServerError, ErrorCodeInternal, "internal server error")
}

// NewValidationErrorFromBinding extracts validation errors from ShouldBindJSON errors
// and returns a descriptive AppError with clean field-level messages.
func NewValidationErrorFromBinding(bindErr error) *AppError {
	// Handle JSON syntax errors
	var syntaxErr *json.SyntaxError
	if errors.As(bindErr, &syntaxErr) {
		return NewValidationError("invalid JSON body")
	}

	// Handle JSON type errors
	var unmarshalTypeErr *json.UnmarshalTypeError
	if errors.As(bindErr, &unmarshalTypeErr) {
		return NewValidationError(
			fmt.Sprintf("invalid value for field '%s'", unmarshalTypeErr.Field),
		)
	}

	var ve validator.ValidationErrors
	if ok := AsValidationErrors(bindErr, &ve); ok && len(ve) > 0 {
		// For single field errors, provide a clean message
		if len(ve) == 1 {
			fe := ve[0]
			field := strings.ToLower(fe.Field())
			switch fe.Tag() {
			case "required":
				return NewValidationError(fmt.Sprintf("%s is required", field))
			case "min":
				if fe.Param() != "" {
					return NewValidationError(fmt.Sprintf("%s must be at least %s", field, fe.Param()))
				}
			case "max":
				if fe.Param() != "" {
					return NewValidationError(fmt.Sprintf("%s must not exceed %s", field, fe.Param()))
				}
			}
		}

		var errMsgs []string
		for _, fe := range ve {
			field := strings.ToLower(fe.Field())
			switch fe.Tag() {
			case "required":
				errMsgs = append(errMsgs, fmt.Sprintf("%s is required", field))
			default:
				errMsgs = append(errMsgs, fmt.Sprintf("field '%s' %s", field, fe.Tag()))
			}
		}
		msg := strings.Join(errMsgs, "; ")
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
