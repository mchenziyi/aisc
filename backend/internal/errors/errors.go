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

// String error codes matching the API spec.
const (
	CodeValidation      = model.ErrCodeValidation
	CodeConflict        = model.ErrCodeConflict
	CodeUnauthorized    = model.ErrCodeUnauthorized
	CodeTokenExpired    = model.ErrCodeTokenExpired
	CodeNotFound        = model.ErrCodeNotFound
	CodeVersionConflict = model.ErrCodeVersionConflict
	CodeInternal        = model.ErrCodeInternal
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
	ErrInternal         = NewAppError(CodeInternal, http.StatusInternalServerError, "服务器内部错误，请稍后重试")
	ErrUnauthorized     = NewAppError(CodeUnauthorized, http.StatusUnauthorized, "未认证或 Token 无效")
	ErrTokenExpired     = NewAppError(CodeTokenExpired, http.StatusUnauthorized, "Token 已过期")
	ErrNotFound         = NewAppError(CodeNotFound, http.StatusNotFound, "资源不存在")
	ErrUsernameTaken    = NewAppError(CodeConflict, http.StatusConflict, "用户名已存在")
	ErrVersionConflict  = NewAppError(CodeVersionConflict, http.StatusConflict, "数据版本冲突，请刷新后重试")
)

// ValidationError creates a validation error with CodeValidation.
func ValidationError(message string) *AppError {
	return NewAppError(CodeValidation, http.StatusBadRequest, message)
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
	if errors.As(bindErr, &ve) && len(ve) > 0 {
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

// ─── Response helpers ─────────────────────────────────────────

// RespondError sends a standardized error response via the Gin context.
// It reads request_id from the Gin context if available.
func RespondError(c *gin.Context, appErr *AppError) {
	requestID := c.GetString("request_id")
	resp := model.ErrorResponse{
		ErrorCode: appErr.Code,
		Message:   appErr.Message,
		RequestID: requestID,
	}
	if len(appErr.Details) > 0 {
		resp.Details = appErr.Details
	}
	c.JSON(appErr.HTTPCode, resp)
}

// RespondSuccess sends a success response with the data directly (no envelope).
// For health check, use RespondHealthCheck instead.
func RespondSuccess(c *gin.Context, httpCode int, data interface{}) {
	c.JSON(httpCode, data)
}
