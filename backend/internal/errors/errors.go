package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"

	"todo-api/internal/model"
)

// Business error codes (matching Tech Design §8.7)
const (
	CodeSuccess           = 0    // 成功
	CodeValidation        = 1001 // 参数校验失败
	CodeUsernameTaken     = 1002 // 用户名已存在
	CodeUnauthorized      = 2001 // 认证失败（Token 缺失/无效/过期）
	CodeRefreshFailed     = 2002 // Refresh Token 无效
	CodeNotFound          = 3001 // 资源不存在（包括越权访问）
	CodeInternal          = 9999 // 服务器内部错误
)

// AppError represents a structured application error.
type AppError struct {
	Code        int                `json:"code"`              // 业务错误码
	Message     string             `json:"message"`           // 错误描述
	HTTPCode    int                `json:"-"`                 // HTTP 状态码（不序列化）
	FieldErrors []model.FieldError `json:"-"`                 // 字段级错误（不序列化到 AppError 自身）
}

func (e *AppError) Error() string {
	return e.Message
}

// NewAppError creates a new AppError with the given business code, HTTP status, and message.
func NewAppError(businessCode, httpCode int, message string) *AppError {
	return &AppError{
		Code:     businessCode,
		HTTPCode: httpCode,
		Message:  message,
	}
}

// Predefined application errors
var (
	ErrInternal       = NewAppError(CodeInternal, http.StatusInternalServerError, "服务器内部错误，请稍后重试")
	ErrUnauthorized   = NewAppError(CodeUnauthorized, http.StatusUnauthorized, "请先登录")
	ErrTokenExpired   = NewAppError(CodeUnauthorized, http.StatusUnauthorized, "Token 已过期，请重新登录")
	ErrInvalidToken   = NewAppError(CodeUnauthorized, http.StatusUnauthorized, "无效的 Token")
	ErrRefreshFailed  = NewAppError(CodeRefreshFailed, http.StatusUnauthorized, "Refresh Token 无效或已过期")
	ErrNotFound       = NewAppError(CodeNotFound, http.StatusNotFound, "资源不存在")
	ErrUsernameTaken  = NewAppError(CodeUsernameTaken, http.StatusConflict, "用户名已存在")
)

// ValidationError creates a validation error with optional field errors.
func ValidationError(message string, fieldErrors ...model.FieldError) *AppError {
	err := NewAppError(CodeValidation, http.StatusBadRequest, message)
	if len(fieldErrors) > 0 {
		err.FieldErrors = fieldErrors
	}
	return err
}

// RespondError sends a standardized error response via the Gin context.
// This is the recommended way to return errors from handlers.
func RespondError(c GinContext, appErr *AppError) {
	resp := model.NewErrorResponse(appErr.Code, appErr.Message, appErr.FieldErrors)
	c.JSON(appErr.HTTPCode, resp)
}

// GinContext defines the minimal interface we need from gin.Context.
type GinContext interface {
	JSON(code int, obj interface{})
}

// NewValidationErrorFromBinding extracts validation errors from ShouldBindJSON errors
// and returns a descriptive AppError with clean field-level messages.
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

		// Create the AppError with field errors
		return ValidationError("请求参数错误", fieldErrors...)
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
