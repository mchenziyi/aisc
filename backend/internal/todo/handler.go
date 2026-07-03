package todo

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	apperrors "todo-api/internal/errors"
	"todo-api/internal/model"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// CreateTodo handles POST /v1/todos
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

	c.JSON(http.StatusCreated, model.ResponseEnvelope{
		Code:    0,
		Data:    resp,
		Message: "ok",
	})
}

// ListTodos handles GET /v1/todos
func (h *Handler) ListTodos(c *gin.Context) {
	userID := c.GetInt64("user_id")

	// Parse and auto-correct page (default 1, minimum 1)
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	// Parse and auto-correct page_size (default 20, min 1, max 100)
	pageSizeStr := c.DefaultQuery("page_size", "20")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 1
	}
	if pageSize > 100 {
		pageSize = 100
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

	c.JSON(http.StatusOK, model.ResponseEnvelope{
		Code:    0,
		Data:    resp,
		Message: "ok",
	})
}

// GetTodo handles GET /v1/todos/:id
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

	c.JSON(http.StatusOK, model.ResponseEnvelope{
		Code:    0,
		Data:    resp,
		Message: "ok",
	})
}

// UpdateTodo handles PUT /v1/todos/:id
func (h *Handler) UpdateTodo(c *gin.Context) {
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

	resp, appErr := h.service.Update(c.Request.Context(), userID, todoID, &req)
	if appErr != nil {
		apperrors.RespondError(c, appErr)
		return
	}

	c.JSON(http.StatusOK, model.ResponseEnvelope{
		Code:    0,
		Data:    resp,
		Message: "ok",
	})
}

// CompleteTodo handles PATCH /v1/todos/:id (mark as completed, idempotent)
func (h *Handler) CompleteTodo(c *gin.Context) {
	userID := c.GetInt64("user_id")

	todoID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || todoID <= 0 {
		apperrors.RespondError(c, apperrors.ValidationError("待办事项 ID 格式无效"))
		return
	}

	resp, appErr := h.service.Complete(c.Request.Context(), userID, todoID)
	if appErr != nil {
		apperrors.RespondError(c, appErr)
		return
	}

	c.JSON(http.StatusOK, model.ResponseEnvelope{
		Code:    0,
		Data:    resp,
		Message: "ok",
	})
}

// DeleteTodo handles DELETE /v1/todos/:id
func (h *Handler) DeleteTodo(c *gin.Context) {
	userID := c.GetInt64("user_id")

	todoID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || todoID <= 0 {
		apperrors.RespondError(c, apperrors.ValidationError("待办事项 ID 格式无效"))
		return
	}

	if appErr := h.service.Delete(c.Request.Context(), userID, todoID); appErr != nil {
		apperrors.RespondError(c, appErr)
		return
	}

	c.Status(http.StatusNoContent)
}
