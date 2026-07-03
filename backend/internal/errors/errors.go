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

// Business error codes matching the tech design spec.
const (
	// 400
	CodeValidation     = 1001
	CodeUsernameTaken  = 1002
	// 401
	CodeUnauthorized   = 2001
	CodeRefreshInvalid = 2002
	CodeTokenExpired   = 2003
	CodeInvalidToken   = 2004
	// 404
	CodeNotFound       = 3001
	// 500
	CodeInternal       = 9999
)

// AppError represents a structured application error.
type AppError struct {
	Code     int               // business error code
	Message  string            // human-readable message
	HTTPCode int               // HTTP status code
	Details  []model.FieldError // optional field-level errors
}

func (e *AppError) Error() string {
	return e.Message
}

// NewAppError creates a new AppError.
func NewAppError(code int, httpCode int, message string) *AppError {
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
)

// ValidationError creates a validation error with code 1001.
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

// RespondError sends a standardized error response via the Gin context.
func RespondError(c *gin.Context, appErr *AppError) {
	resp := &model.ErrorResponse{
		Code:    appErr.Code,
		Message: appErr.Message,
	}
	if len(appErr.Details) > 0 {
		resp.Errors = appErr.Details
	}
	c.JSON(appErr.HTTPCode, resp)
}
