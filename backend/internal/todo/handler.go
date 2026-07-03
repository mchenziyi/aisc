package todo

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	apperrors "todo-api/internal/errors"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// CreateTodo handles POST /api/v1/todos
func (h *Handler) CreateTodo(c *gin.Context) {
	userID := c.GetInt64("user_id")

	var req CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.RespondError(c, apperrors.NewValidationErrorFromBinding(err))
		return
	}

	resp, appErr := h.service.Create(c.Request.Context(), userID, &req)
	if appErr != nil {
		apperrors.RespondError(c, appErr)
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// ListTodos handles GET /api/v1/todos
func (h *Handler) ListTodos(c *gin.Context) {
	userID := c.GetInt64("user_id")

	// Parse page (default 1, must be >= 1)
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		apperrors.RespondError(c, apperrors.ValidationError("page 参数无效，必须为正整数"))
		return
	}

	// Parse page_size (default 20, min 1, max 100)
	pageSizeStr := c.DefaultQuery("page_size", "20")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		apperrors.RespondError(c, apperrors.ValidationError("page_size 参数无效，必须为正整数"))
		return
	}
	if pageSize > 100 {
		apperrors.RespondError(c, apperrors.ValidationError("page_size 不能超过 100"))
		return
	}

	// Parse status filter (default "all")
	status := c.DefaultQuery("status", "all")
	if status != "all" && status != "pending" && status != "completed" {
		apperrors.RespondError(c, apperrors.ValidationError("status 参数无效，只能是 all、pending 或 completed"))
		return
	}

	q := &ListQuery{
		Page:     page,
		PageSize: pageSize,
		Status:   status,
	}

	resp, appErr := h.service.List(c.Request.Context(), userID, q)
	if appErr != nil {
		apperrors.RespondError(c, appErr)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetTodo handles GET /api/v1/todos/:id
func (h *Handler) GetTodo(c *gin.Context) {
	userID := c.GetInt64("user_id")

	todoID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || todoID <= 0 {
		apperrors.RespondError(c, apperrors.ValidationError("待办事项 ID 格式无效"))
		return
	}

	resp, appErr := h.service.GetByID(c.Request.Context(), userID, todoID)
	if appErr != nil {
		apperrors.RespondError(c, appErr)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// PatchTodo handles PATCH /api/v1/todos/:id (update with optimistic locking)
func (h *Handler) PatchTodo(c *gin.Context) {
	userID := c.GetInt64("user_id")

	todoID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || todoID <= 0 {
		apperrors.RespondError(c, apperrors.ValidationError("待办事项 ID 格式无效"))
		return
	}

	var req UpdateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.RespondError(c, apperrors.NewValidationErrorFromBinding(err))
		return
	}

	resp, appErr := h.service.Patch(c.Request.Context(), userID, todoID, &req)
	if appErr != nil {
		apperrors.RespondError(c, appErr)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteTodo handles DELETE /api/v1/todos/:id (with optimistic locking via version query param)
func (h *Handler) DeleteTodo(c *gin.Context) {
	userID := c.GetInt64("user_id")

	todoID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || todoID <= 0 {
		apperrors.RespondError(c, apperrors.ValidationError("待办事项 ID 格式无效"))
		return
	}

	// Parse version query parameter for optimistic locking
	versionStr := c.Query("version")
	if versionStr == "" {
		apperrors.RespondError(c, apperrors.ValidationError("version 参数不能为空"))
		return
	}
	version, err := strconv.Atoi(versionStr)
	if err != nil || version < 1 {
		apperrors.RespondError(c, apperrors.ValidationError("version 参数无效，必须为正整数"))
		return
	}

	if appErr := h.service.Delete(c.Request.Context(), userID, todoID, version); appErr != nil {
		apperrors.RespondError(c, appErr)
		return
	}

	c.Status(http.StatusNoContent)
}
