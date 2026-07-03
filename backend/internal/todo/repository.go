package todo

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// Create inserts a new todo record and returns the created todo.
func (r *Repository) Create(ctx context.Context, todo *Todo) (*Todo, error) {
	var t Todo
	err := r.pool.QueryRow(ctx,
		`INSERT INTO todos (user_id, title, description, due_date)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, user_id, title, description, due_date, is_completed, completed_at, created_at, updated_at`,
		todo.UserID, todo.Title, todo.Description, todo.DueDate,
	).Scan(&t.ID, &t.UserID, &t.Title, &t.Description, &t.DueDate, &t.IsCompleted, &t.CompletedAt, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// FindByIDAndUser finds a todo by id and user_id (ensures data isolation).
func (r *Repository) FindByIDAndUser(ctx context.Context, id, userID int64) (*Todo, error) {
	var t Todo
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, title, description, due_date, is_completed, completed_at, created_at, updated_at
		 FROM todos WHERE id = $1 AND user_id = $2`,
		id, userID,
	).Scan(&t.ID, &t.UserID, &t.Title, &t.Description, &t.DueDate, &t.IsCompleted, &t.CompletedAt, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

// ListByUser returns paginated todos for a user, ordered by created_at DESC.
// status: "all", "pending", "completed".
func (r *Repository) ListByUser(ctx context.Context, userID int64, q *ListQuery) ([]*Todo, int64, error) {
	// Build count query
	countWhere := "WHERE user_id = $1"
	countArgs := []interface{}{userID}

	// Build list query
	listWhere := "WHERE user_id = $1"
	listArgs := []interface{}{userID}

	argIdx := 2

	if q.Status == "pending" {
		countWhere += " AND is_completed = FALSE"
		listWhere += " AND is_completed = FALSE"
	} else if q.Status == "completed" {
		countWhere += " AND is_completed = TRUE"
		listWhere += " AND is_completed = TRUE"
	}

	var total int64
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM todos `+countWhere,
		countArgs...,
	).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []*Todo{}, 0, nil
	}

	offset := (q.Page - 1) * q.PageSize

	query := `SELECT id, user_id, title, description, due_date, is_completed, completed_at, created_at, updated_at
		 FROM todos ` + listWhere + `
		 ORDER BY created_at DESC, id DESC
		 LIMIT $` + itoa(argIdx) + ` OFFSET $` + itoa(argIdx+1)

	listArgs = append(listArgs, q.PageSize, offset)

	rows, err := r.pool.Query(ctx, query, listArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var todos []*Todo
	for rows.Next() {
		var t Todo
		if err := rows.Scan(&t.ID, &t.UserID, &t.Title, &t.Description, &t.DueDate, &t.IsCompleted, &t.CompletedAt, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, 0, err
		}
		todos = append(todos, &t)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return todos, total, nil
}

// Update updates a todo's fields. Only non-nil fields will be updated.
func (r *Repository) Update(ctx context.Context, todoID, userID int64, title *string, description *string, dueDate *time.Time) (*Todo, error) {
	// Build dynamic SET clause
	setClauses := []string{}
	args := []interface{}{}
	argIdx := 1

	if title != nil {
		setClauses = append(setClauses, "title = $"+itoa(argIdx))
		args = append(args, *title)
		argIdx++
	}
	if description != nil {
		setClauses = append(setClauses, "description = $"+itoa(argIdx))
		if *description == "" {
			args = append(args, nil)
		} else {
			args = append(args, *description)
		}
		argIdx++
	}
	if dueDate != nil {
		if dueDate.IsZero() {
			setClauses = append(setClauses, "due_date = NULL")
		} else {
			setClauses = append(setClauses, "due_date = $"+itoa(argIdx))
			args = append(args, *dueDate)
			argIdx++
		}
	}

	if len(setClauses) == 0 {
		// Nothing to update; fetch current state
		return r.FindByIDAndUser(ctx, todoID, userID)
	}

	idParam := argIdx
	args = append(args, todoID)
	argIdx++

	userIDParam := argIdx
	args = append(args, userID)

	sql := `UPDATE todos SET ` + joinStrings(setClauses, ", ") +
		` WHERE id = $` + itoa(idParam) +
		` AND user_id = $` + itoa(userIDParam) +
		` RETURNING id, user_id, title, description, due_date, is_completed, completed_at, created_at, updated_at`

	var t Todo
	err := r.pool.QueryRow(ctx, sql, args...).Scan(
		&t.ID, &t.UserID, &t.Title, &t.Description, &t.DueDate, &t.IsCompleted, &t.CompletedAt, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

// Complete marks a todo as completed (idempotent).
// Uses COALESCE to ensure completed_at is only set on first completion.
func (r *Repository) Complete(ctx context.Context, todoID, userID int64) (*Todo, error) {
	var t Todo
	err := r.pool.QueryRow(ctx,
		`UPDATE todos
		 SET is_completed = TRUE,
		     completed_at = COALESCE(completed_at, NOW())
		 WHERE id = $1 AND user_id = $2
		 RETURNING id, user_id, title, description, due_date, is_completed, completed_at, created_at, updated_at`,
		todoID, userID,
	).Scan(&t.ID, &t.UserID, &t.Title, &t.Description, &t.DueDate, &t.IsCompleted, &t.CompletedAt, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

// Delete deletes a todo by id and user_id.
// Returns false if no rows were deleted.
func (r *Repository) Delete(ctx context.Context, id, userID int64) (bool, error) {
	ct, err := r.pool.Exec(ctx,
		`DELETE FROM todos WHERE id = $1 AND user_id = $2`,
		id, userID,
	)
	if err != nil {
		return false, err
	}
	return ct.RowsAffected() == 1, nil
}

// ─── Helper functions ────────────────────────────────────────

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	var buf [12]byte
	pos := len(buf)
	neg := i < 0
	if neg {
		i = -i
	}
	for i > 0 {
		pos--
		buf[pos] = byte('0' + i%10)
		i /= 10
	}
	if neg {
		pos--
		buf[pos] = '-'
	}
	return string(buf[pos:])
}

func joinStrings(elems []string, sep string) string {
	if len(elems) == 0 {
		return ""
	}
	if len(elems) == 1 {
		return elems[0]
	}
	n := len(sep) * (len(elems) - 1)
	for _, e := range elems {
		n += len(e)
	}
	b := make([]byte, 0, n)
	b = append(b, elems[0]...)
	for _, e := range elems[1:] {
		b = append(b, sep...)
		b = append(b, e...)
	}
	return string(b)
}
