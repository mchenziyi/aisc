package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	apperrors "todo-api/internal/errors"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Register handles POST /api/v1/users
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.RespondError(c, apperrors.NewValidationErrorFromBinding(err))
		return
	}

	resp, appErr := h.service.Register(c.Request.Context(), &req)
	if appErr != nil {
		apperrors.RespondError(c, appErr)
		return
	}

	apperrors.RespondSuccess(c, http.StatusCreated, resp)
}

// Login handles POST /api/v1/auth/login
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.RespondError(c, apperrors.NewValidationErrorFromBinding(err))
		return
	}

	resp, appErr := h.service.Login(c.Request.Context(), &req)
	if appErr != nil {
		apperrors.RespondError(c, appErr)
		return
	}

	apperrors.RespondSuccess(c, http.StatusOK, resp)
}

// RefreshToken handles POST /api/v1/auth/refresh
func (h *Handler) RefreshToken(c *gin.Context) {
	userID := c.GetInt64("user_id")

	resp, appErr := h.service.RefreshToken(c.Request.Context(), userID)
	if appErr != nil {
		apperrors.RespondError(c, appErr)
		return
	}

	apperrors.RespondSuccess(c, http.StatusOK, resp)
}

// GetCurrentUser handles GET /api/v1/users/me
func (h *Handler) GetCurrentUser(c *gin.Context) {
	userID := c.GetInt64("user_id")

	user, appErr := h.service.GetCurrentUser(c.Request.Context(), userID)
	if appErr != nil {
		apperrors.RespondError(c, appErr)
		return
	}

	apperrors.RespondSuccess(c, http.StatusOK, user)
}
