package todo

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	apperrors "todo-api/internal/errors"
	"todo-api/internal/model"
)

// ISO 8601 date-time pattern (simplified)
var dateTimeRegex = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}(T\d{2}:\d{2}:\d{2}(\.\d+)?(Z|[+-]\d{2}:\d{2})?)?$`)

const (
	maxTitleLen       = 200
	maxDescriptionLen = 500
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// Create creates a new todo for the given user.
func (s *Service) Create(ctx context.Context, userID int64, req *CreateTodoRequest) (*model.TodoResponse, *apperrors.AppError) {
	// Validate title
	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, apperrors.ValidationError("标题不能为空")
	}
	if utf8.RuneCountInString(title) > maxTitleLen {
		return nil, apperrors.ValidationError("标题不能超过 200 个字符")
	}

	// Validate description
	var desc *string
	if req.Description != nil {
		trimmed := strings.TrimSpace(*req.Description)
		if trimmed != "" {
			if utf8.RuneCountInString(trimmed) > maxDescriptionLen {
				return nil, apperrors.ValidationError("描述不能超过 500 个字符")
			}
			desc = &trimmed
		}
	}

	// Validate and parse due_date if provided
	var dueDate *time.Time
	if req.DueDate != nil {
		parsed, err := parseDueDate(*req.DueDate)
		if err != nil {
			return nil, apperrors.ValidationError(err.Error())
		}
		dueDate = parsed
	}

	todo := &Todo{
		UserID:      userID,
		Title:       title,
		Description: desc,
		DueDate:     dueDate,
		IsCompleted: false,
		CompletedAt: nil,
	}

	created, err := s.repo.Create(ctx, todo)
	if err != nil {
		log.Printf("internal error: failed to create todo for user %d: %v", userID, err)
		return nil, apperrors.ErrInternal
	}

	return todoToResponse(created), nil
}

// GetByID retrieves a single todo by ID, ensuring it belongs to the user.
func (s *Service) GetByID(ctx context.Context, userID, todoID int64) (*model.TodoResponse, *apperrors.AppError) {
	todo, err := s.repo.FindByIDAndUser(ctx, todoID, userID)
	if err != nil {
		log.Printf("internal error: failed to find todo %d for user %d: %v", todoID, userID, err)
		return nil, apperrors.ErrInternal
	}
	if todo == nil {
		return nil, apperrors.ErrNotFound
	}
	return todoToResponse(todo), nil
}

// List returns a paginated list of todos for the user, with optional status filter.
func (s *Service) List(ctx context.Context, userID int64, page, pageSize int, status string) (*model.PaginatedTodos, *apperrors.AppError) {
	// Validate and normalize status
	status = strings.TrimSpace(strings.ToLower(status))
	switch status {
	case "all", "pending", "completed", "":
		// valid
	default:
		return nil, apperrors.ValidationError("状态参数无效，可选值: all, pending, completed")
	}
	if status == "" || status == "all" {
		status = "all"
	}

	todos, total, err := s.repo.ListByUser(ctx, userID, page, pageSize, status)
	if err != nil {
		log.Printf("internal error: failed to list todos for user %d: %v", userID, err)
		return nil, apperrors.ErrInternal
	}

	totalPages := 0
	if total > 0 {
		totalPages = int((total + int64(pageSize) - 1) / int64(pageSize))
	}

	items := make([]model.TodoResponse, 0, len(todos))
	for _, t := range todos {
		items = append(items, *todoToResponse(t))
	}

	return &model.PaginatedTodos{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// Update updates a todo's fields (title, description, due_date).
func (s *Service) Update(ctx context.Context, userID, todoID int64, req *UpdateTodoRequest) (*model.TodoResponse, *apperrors.AppError) {
	// Validate that at least one field is provided
	if req.Title == nil && req.Description == nil && req.DueDate == nil {
		return nil, apperrors.ValidationError("至少需要提供一个要更新的字段")
	}

	// Validate title if provided
	var title *string
	if req.Title != nil {
		t := strings.TrimSpace(*req.Title)
		if t == "" {
			return nil, apperrors.ValidationError("标题不能为空")
		}
		if utf8.RuneCountInString(t) > maxTitleLen {
			return nil, apperrors.ValidationError("标题不能超过 200 个字符")
		}
		title = &t
	}

	// Validate description if provided
	var desc *string
	if req.Description != nil {
		d := strings.TrimSpace(*req.Description)
		if utf8.RuneCountInString(d) > maxDescriptionLen {
			return nil, apperrors.ValidationError("描述不能超过 500 个字符")
		}
		desc = &d
	}

	// Validate and parse due_date if provided
	var dueDate *time.Time
	if req.DueDate != nil {
		if *req.DueDate == "" {
			// Empty string means clear the due_date - store a zero time to signal NULL
			zeroTime := time.Time{}
			dueDate = &zeroTime
		} else {
			parsed, err := parseDueDate(*req.DueDate)
			if err != nil {
				return nil, apperrors.ValidationError(err.Error())
			}
			dueDate = parsed
		}
	}

	updated, err := s.repo.Update(ctx, todoID, userID, title, desc, dueDate)
	if err != nil {
		log.Printf("internal error: failed to update todo %d for user %d: %v", todoID, userID, err)
		return nil, apperrors.ErrInternal
	}
	if updated == nil {
		return nil, apperrors.ErrNotFound
	}

	return todoToResponse(updated), nil
}

// Complete marks a todo as completed. Idempotent.
func (s *Service) Complete(ctx context.Context, userID, todoID int64) (*model.TodoResponse, *apperrors.AppError) {
	updated, err := s.repo.Complete(ctx, todoID, userID)
	if err != nil {
		log.Printf("internal error: failed to complete todo %d for user %d: %v", todoID, userID, err)
		return nil, apperrors.ErrInternal
	}
	if updated == nil {
		return nil, apperrors.ErrNotFound
	}
	return todoToResponse(updated), nil
}

// Delete deletes a todo.
func (s *Service) Delete(ctx context.Context, userID, todoID int64) *apperrors.AppError {
	deleted, err := s.repo.Delete(ctx, todoID, userID)
	if err != nil {
		log.Printf("internal error: failed to delete todo %d for user %d: %v", todoID, userID, err)
		return apperrors.ErrInternal
	}
	if !deleted {
		return apperrors.ErrNotFound
	}
	return nil
}

// parseDueDate validates and parses a due date string.
// Returns nil if the input is empty.
// Returns an error if the format is invalid or the date is in the past.
func parseDueDate(s string) (*time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, nil
	}

	if !dateTimeRegex.MatchString(s) {
		return nil, fmt.Errorf("截止日期格式无效，请使用 ISO 8601 格式")
	}

	// Try parsing with various ISO 8601 formats
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02T15:04:05Z",
		"2006-01-02",
	}

	var parsedTime time.Time
	var err error
	for _, format := range formats {
		parsedTime, err = time.Parse(format, s)
		if err == nil {
			break
		}
	}
	if err != nil {
		return nil, fmt.Errorf("截止日期格式无效，请使用 ISO 8601 格式")
	}

	if parsedTime.Before(time.Now()) {
		return nil, fmt.Errorf("截止日期必须是未来时间")
	}

	return &parsedTime, nil
}

// todoToResponse converts a Todo model to a model.TodoResponse.
func todoToResponse(t *Todo) *model.TodoResponse {
	resp := &model.TodoResponse{
		ID:          t.ID,
		Title:       t.Title,
		IsCompleted: t.IsCompleted,
		UserID:      t.UserID,
		CreatedAt:   t.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   t.UpdatedAt.Format(time.RFC3339),
	}

	if t.Description != nil {
		resp.Description = t.Description
	}
	if t.DueDate != nil {
		s := t.DueDate.Format(time.RFC3339)
		resp.DueDate = &s
	}
	if t.CompletedAt != nil {
		s := t.CompletedAt.Format(time.RFC3339)
		resp.CompletedAt = &s
	}

	return resp
}
