package todo

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"
	"unicode/utf8"

	apperrors "todo-api/internal/errors"
	"todo-api/internal/model"
)

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
	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, apperrors.ValidationError("标题不能为空")
	}
	if utf8.RuneCountInString(title) > maxTitleLen {
		return nil, apperrors.ValidationError("标题不能超过 200 个字符")
	}

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

	var dueDate *time.Time
	if req.DueDate != nil {
		parsed, err := parseDateTime(*req.DueDate)
		if err != nil {
			return nil, apperrors.ValidationError(err.Error())
		}
		if parsed != nil && parsed.Before(time.Now()) {
			return nil, apperrors.ValidationError("截止日期必须是未来时间")
		}
		dueDate = parsed
	}

	todo := &Todo{
		UserID:      userID,
		Title:       title,
		Description: desc,
		DueDate:     dueDate,
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

// List returns a paginated list of todos for the user, with optional status filtering.
func (s *Service) List(ctx context.Context, userID int64, q *ListQuery) (*model.PaginatedTodos, *apperrors.AppError) {
	todos, total, err := s.repo.ListByUser(ctx, userID, q)
	if err != nil {
		log.Printf("internal error: failed to list todos for user %d: %v", userID, err)
		return nil, apperrors.ErrInternal
	}

	totalPages := 0
	if total > 0 {
		totalPages = int((total + int64(q.PageSize) - 1) / int64(q.PageSize))
	}

	items := make([]model.TodoResponse, 0, len(todos))
	for _, t := range todos {
		items = append(items, *todoToResponse(t))
	}

	return &model.PaginatedTodos{
		Items:      items,
		Total:      total,
		Page:       q.Page,
		PageSize:   q.PageSize,
		TotalPages: totalPages,
	}, nil
}

// Update updates a todo's title, description, and/or due_date (PUT).
func (s *Service) Update(ctx context.Context, userID, todoID int64, req *UpdateTodoRequest) (*model.TodoResponse, *apperrors.AppError) {
	// At least one field must be provided
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
		if d == "" {
			// Empty string means clear the description (set to NULL)
			empty := ""
			desc = &empty
		} else {
			if utf8.RuneCountInString(d) > maxDescriptionLen {
				return nil, apperrors.ValidationError("描述不能超过 500 个字符")
			}
			desc = &d
		}
	}

	// Validate and parse due_date if provided
	var dueDate *time.Time
	if req.DueDate != nil {
		dd := strings.TrimSpace(*req.DueDate)
		if dd == "" {
			// Empty string means clear the due_date (set to NULL)
			zero := time.Time{}
			dueDate = &zero
		} else {
			parsed, err := parseDateTime(dd)
			if err != nil {
				return nil, apperrors.ValidationError(err.Error())
			}
			if parsed != nil && parsed.Before(time.Now()) {
				return nil, apperrors.ValidationError("截止日期必须是未来时间")
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

// Complete marks a todo as completed (idempotent).
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

// Delete deletes a todo by id and user_id.
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

// parseDateTime validates and parses a date/time string (ISO 8601).
func parseDateTime(s string) (*time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, nil
	}

	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05",
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
		return nil, fmt.Errorf("截止日期格式无效，请使用 ISO 8601 格式 (如 2025-12-31T23:59:59Z)")
	}

	return &parsedTime, nil
}

// todoToResponse converts a Todo model to a model.TodoResponse matching the API spec.
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
