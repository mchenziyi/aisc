package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Handler handles auth HTTP requests.
type Handler struct {
	service *Service
}

// NewHandler creates a new auth handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Register handles POST /auth/register
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	user, err := h.service.Register(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, ErrUsernameTaken) {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{
				"code":    409,
				"message": "username already exists",
			})
			return
		}
		// Validation errors (like invalid username or password format)
		if isValidationError(err) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": err.Error(),
			})
			return
		}
		// Unexpected server errors
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "internal server error",
		})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// Login handles POST /auth/login
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	resp, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		// Invalid credentials is a client error
		if strings.Contains(err.Error(), "invalid username or password") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "invalid username or password",
			})
			return
		}
		// Unexpected server errors
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Me handles GET /auth/me
func (h *Handler) Me(c *gin.Context) {
	userID := c.GetInt64("user_id")

	user, err := h.service.GetMe(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "user not found",
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

// isValidationError checks if an error is a client-side validation error.
func isValidationError(err error) bool {
	msg := err.Error()
	return strings.Contains(msg, "must be") ||
		strings.Contains(msg, "is required") ||
		strings.Contains(msg, "cannot be")
}
