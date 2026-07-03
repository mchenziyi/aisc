package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"todo-api/internal/model"
)

// Business error codes (string) matching the tech design spec.
const (
	// 400
	CodeValidation     = "VALIDATION_ERROR"
	CodeUsernameTaken  = "USERNAME_TAKEN"
	CodeInvalidParams  = "INVALID_PARAMS"
	// 401
	CodeUnauthorized   = "UNAUTHORIZED"
	CodeRefreshInvalid = "REFRESH_INVALID"
	CodeTokenExpired   = "TOKEN_EXPIRED"
	CodeInvalidToken   = "INVALID_TOKEN"
	// 404
	CodeNotFound       = "NOT_FOUND"
	// 409
	CodeConflict       = "CONFLICT"
	// 500
	CodeInternal       = "INTERNAL_ERROR"
)

// AppError represents a structured application error.
type AppError struct {
	Code     string             // business error code (string)
	Message  string             // human-readable message
	HTTPCode int                // HTTP status code
	Details  []model.FieldError // optional field-level errors
}

func (e *AppError) Error() string {
	return e.Message
}

// NewAppError creates a new AppError.
func NewAppError(code string, httpCode int, message string) *AppError {
	return &AppError{
		Code:     code,
		HTTPCode: httpCode,
		Message:  message,
	}
}

// Predefined application errors.
var (
	ErrInternal       = NewAppError(CodeInternal, http.StatusInternalServerError, "服务器内部错误，请稍后重试")
	ErrUnauthorized   = NewAppError(CodeUnauthorized, http.StatusUnauthorized, "请先登录")
	ErrTokenExpired   = NewAppError(CodeTokenExpired, http.StatusUnauthorized, "Token 已过期，请重新登录")
	ErrInvalidToken   = NewAppError(CodeInvalidToken, http.StatusUnauthorized, "无效的 Token")
	ErrRefreshInvalid = NewAppError(CodeRefreshInvalid, http.StatusUnauthorized, "Refresh Token 无效或已过期")
	ErrNotFound       = NewAppError(CodeNotFound, http.StatusNotFound, "资源不存在")
	ErrUsernameTaken  = NewAppError(CodeUsernameTaken, http.StatusConflict, "用户名已存在")
	ErrConflict       = NewAppError(CodeConflict, http.StatusConflict, "数据冲突，请刷新后重试")
)

// ValidationError creates a validation error with code VALIDATION_ERROR.
func ValidationError(message string) *AppError {
	return NewAppError(CodeValidation, http.StatusBadRequest, message)
}

// ConflictError creates a conflict error with code CONFLICT.
func ConflictError(message string) *AppError {
	return NewAppError(CodeConflict, http.StatusConflict, message)
}

// IsConflictError checks if the given error is a conflict error.
func IsConflictError(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == CodeConflict
	}
	return false
}

// ValidationErrorWithFields creates a validation error with field-level details.
func ValidationErrorWithFields(message string, fields []model.FieldError) *AppError {
	return &AppError{
		Code:     CodeValidation,
		HTTPCode: http.StatusBadRequest,
		Message:  message,
		Details:  fields,
	}
}

// NewValidationErrorFromBinding extracts validation errors from ShouldBindJSON errors.
func NewValidationErrorFromBinding(bindErr error) *AppError {
	// Handle JSON syntax errors
	var syntaxErr *json.SyntaxError
	if errors.As(bindErr, &syntaxErr) {
		return ValidationError("请求体格式错误")
	}

	// Handle JSON type errors
	var unmarshalTypeErr *json.UnmarshalTypeError
	if errors.As(bindErr, &unmarshalTypeErr) {
		return ValidationError(
			fmt.Sprintf("字段 '%s' 的值类型不正确", unmarshalTypeErr.Field),
		)
	}

	// Handle validator errors
	var ve validator.ValidationErrors
	if ok := AsValidationErrors(bindErr, &ve); ok && len(ve) > 0 {
		var fieldErrors []model.FieldError
		for _, fe := range ve {
			field := strings.ToLower(fe.Field())
			var msg string
			switch fe.Tag() {
			case "required":
				msg = fmt.Sprintf("%s 不能为空", field)
			case "min":
				msg = fmt.Sprintf("%s 长度不能少于 %s", field, fe.Param())
			case "max":
				msg = fmt.Sprintf("%s 长度不能超过 %s", field, fe.Param())
			case "oneof":
				msg = fmt.Sprintf("%s 只能是 %s", field, fe.Param())
			case "len":
				msg = fmt.Sprintf("%s 长度必须为 %s", field, fe.Param())
			default:
				msg = fmt.Sprintf("字段 '%s' 校验失败 (%s)", field, fe.Tag())
			}
			fieldErrors = append(fieldErrors, model.FieldError{Field: field, Message: msg})
		}
		return ValidationErrorWithFields("请求参数错误", fieldErrors)
	}

	return ValidationError("无效的请求体")
}

// AsValidationErrors checks if the error is a validator.ValidationErrors.
var AsValidationErrors = func(err error, target *validator.ValidationErrors) bool {
	if err == nil {
		return false
	}
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		*target = ve
		return true
	}
	return false
}

// GetRequestID retrieves the request_id from the Gin context.
func GetRequestID(c *gin.Context) string {
	if id, exists := c.Get("request_id"); exists {
		if s, ok := id.(string); ok {
			return s
		}
	}
	return ""
}

// RespondError sends a standardized error response via the Gin context.
func RespondError(c *gin.Context, appErr *AppError) {
	resp := &model.ErrorResponse{
		ErrorCode: appErr.Code,
		Message:   appErr.Message,
		RequestID: GetRequestID(c),
	}
	if len(appErr.Details) > 0 {
		resp.Details = appErr.Details
	}
	c.JSON(appErr.HTTPCode, resp)
}
