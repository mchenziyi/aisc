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
		c.Error(apperrors.NewValidationErrorFromBinding(err))
		c.Abort()
		return
	}

	resp, appErr := h.service.Create(c.Request.Context(), userID, &req)
	if appErr != nil {
		_ = c.Error(appErr)
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, model.NewSuccessResponse(resp))
}

// ListTodos handles GET /v1/todos
func (h *Handler) ListTodos(c *gin.Context) {
	userID := c.GetInt64("user_id")

	// Parse and normalize page
	page, err := parseQueryInt(c, "page", 1)
	if err != nil || page < 1 {
		page = 1
	}

	// Parse and normalize page_size
	pageSize, err := parseQueryInt(c, "page_size", 20)
	if err != nil || pageSize < 1 {
		pageSize = 1
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// Read status filter
	status := c.DefaultQuery("status", "all")

	resp, appErr := h.service.List(c.Request.Context(), userID, page, pageSize, status)
	if appErr != nil {
		_ = c.Error(appErr)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(resp))
}

// GetTodo handles GET /v1/todos/:id
func (h *Handler) GetTodo(c *gin.Context) {
	userID := c.GetInt64("user_id")

	todoID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(apperrors.ValidationError("待办事项 ID 格式无效"))
		c.Abort()
		return
	}

	resp, appErr := h.service.GetByID(c.Request.Context(), userID, todoID)
	if appErr != nil {
		_ = c.Error(appErr)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(resp))
}

// UpdateTodo handles PUT /v1/todos/:id
func (h *Handler) UpdateTodo(c *gin.Context) {
	userID := c.GetInt64("user_id")

	todoID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(apperrors.ValidationError("待办事项 ID 格式无效"))
		c.Abort()
		return
	}

	var req UpdateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperrors.NewValidationErrorFromBinding(err))
		c.Abort()
		return
	}

	resp, appErr := h.service.Update(c.Request.Context(), userID, todoID, &req)
	if appErr != nil {
		_ = c.Error(appErr)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(resp))
}

// CompleteTodo handles PATCH /v1/todos/:id
func (h *Handler) CompleteTodo(c *gin.Context) {
	userID := c.GetInt64("user_id")

	todoID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(apperrors.ValidationError("待办事项 ID 格式无效"))
		c.Abort()
		return
	}

	resp, appErr := h.service.Complete(c.Request.Context(), userID, todoID)
	if appErr != nil {
		_ = c.Error(appErr)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(resp))
}

// DeleteTodo handles DELETE /v1/todos/:id
func (h *Handler) DeleteTodo(c *gin.Context) {
	userID := c.GetInt64("user_id")

	todoID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(apperrors.ValidationError("待办事项 ID 格式无效"))
		c.Abort()
		return
	}

	appErr := h.service.Delete(c.Request.Context(), userID, todoID)
	if appErr != nil {
		_ = c.Error(appErr)
		c.Abort()
		return
	}

	c.Status(http.StatusNoContent)
}

// parseQueryInt parses an integer query parameter with a default value.
func parseQueryInt(c *gin.Context, key string, defaultVal int) (int, error) {
	val := c.DefaultQuery(key, "")
	if val == "" {
		return defaultVal, nil
	}
	return strconv.Atoi(val)
}
