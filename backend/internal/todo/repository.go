package todo

import (
	"context"
	"errors"
	"fmt"
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
		 RETURNING id, user_id, title, description, due_date, is_completed AS completed, completed_at, version, created_at, updated_at`,
		todo.UserID, todo.Title, todo.Description, todo.DueDate,
	).Scan(&t.ID, &t.UserID, &t.Title, &t.Description, &t.DueDate, &t.Completed, &t.CompletedAt, &t.Version, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// FindByIDAndUser finds a todo by id and user_id (ensures data isolation).
func (r *Repository) FindByIDAndUser(ctx context.Context, id, userID int64) (*Todo, error) {
	var t Todo
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, title, description, due_date, is_completed AS completed, completed_at, version, created_at, updated_at
		 FROM todos WHERE id = $1 AND user_id = $2`,
		id, userID,
	).Scan(&t.ID, &t.UserID, &t.Title, &t.Description, &t.DueDate, &t.Completed, &t.CompletedAt, &t.Version, &t.CreatedAt, &t.UpdatedAt)
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
	countWhere := "WHERE user_id = $1"
	countArgs := []interface{}{userID}

	listWhere := "WHERE user_id = $1"

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

	query := `SELECT id, user_id, title, description, due_date, is_completed AS completed, completed_at, version, created_at, updated_at
		 FROM todos ` + listWhere + `
		 ORDER BY created_at DESC, id DESC
		 LIMIT $2 OFFSET $3`

	rows, err := r.pool.Query(ctx, query, userID, q.PageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var todos []*Todo
	for rows.Next() {
		var t Todo
		if err := rows.Scan(&t.ID, &t.UserID, &t.Title, &t.Description, &t.DueDate, &t.Completed, &t.CompletedAt, &t.Version, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, 0, err
		}
		todos = append(todos, &t)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return todos, total, nil
}

// Patch updates a todo's fields (title, description, due_date, completed) with optimistic locking.
// Returns the updated todo, or nil if not found.
// Returns an error with "version_conflict" if the version doesn't match.
func (r *Repository) Patch(ctx context.Context, todoID, userID int64, version int, title *string, description *string, dueDate *time.Time, completed *bool) (*Todo, error) {
	// Build dynamic SET clause
	setClauses := []string{}
	args := []interface{}{}
	argIdx := 1

	args = append(args, todoID)
	argIdx++ // $1 = todoID

	args = append(args, userID)
	argIdx++ // $2 = userID

	if title != nil {
		setClauses = append(setClauses, fmt.Sprintf("title = $%d", argIdx))
		args = append(args, *title)
		argIdx++
	}
	if description != nil {
		if *description == "" {
			setClauses = append(setClauses, "description = NULL")
		} else {
			setClauses = append(setClauses, fmt.Sprintf("description = $%d", argIdx))
			args = append(args, *description)
			argIdx++
		}
	}
	if dueDate != nil {
		if dueDate.IsZero() {
			setClauses = append(setClauses, "due_date = NULL")
		} else {
			setClauses = append(setClauses, fmt.Sprintf("due_date = $%d", argIdx))
			args = append(args, *dueDate)
			argIdx++
		}
	}
	if completed != nil {
		setClauses = append(setClauses, fmt.Sprintf("is_completed = $%d", argIdx))
		args = append(args, *completed)
		argIdx++
		if *completed {
			setClauses = append(setClauses, "completed_at = COALESCE(completed_at, NOW())")
		}
	}

	if len(setClauses) == 0 {
		// Nothing to update, return existing todo
		return r.FindByIDAndUser(ctx, todoID, userID)
	}

	// Add version increment
	setClauses = append(setClauses, fmt.Sprintf("version = version + 1"))

	// Version check
	args = append(args, version)
	versionParam := argIdx
	argIdx++

	sql := `UPDATE todos SET ` + joinStrings(setClauses, ", ") +
		` WHERE id = $1 AND user_id = $2 AND version = $` + itoa(versionParam) +
		` RETURNING id, user_id, title, description, due_date, is_completed AS completed, completed_at, version, created_at, updated_at`

	var t Todo
	err := r.pool.QueryRow(ctx, sql, args...).Scan(
		&t.ID, &t.UserID, &t.Title, &t.Description, &t.DueDate, &t.Completed, &t.CompletedAt, &t.Version, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Check if the row exists at all (version mismatch vs not found)
			exists, checkErr := r.existsByIDAndUser(ctx, todoID, userID)
			if checkErr != nil {
				return nil, checkErr
			}
			if !exists {
				return nil, nil
			}
			// Row exists but version mismatch → conflict
			return nil, fmt.Errorf("version_conflict")
		}
		return nil, err
	}
	return &t, nil
}

// Delete deletes a todo by id and user_id with optimistic locking.
// Returns true if deleted, false if not found.
// Returns an error with "version_conflict" if the version doesn't match.
func (r *Repository) Delete(ctx context.Context, id, userID int64, version int) (bool, error) {
	ct, err := r.pool.Exec(ctx,
		`DELETE FROM todos WHERE id = $1 AND user_id = $2 AND version = $3`,
		id, userID, version,
	)
	if err != nil {
		return false, err
	}
	if ct.RowsAffected() == 0 {
		// Check if the row exists at all
		exists, checkErr := r.existsByIDAndUser(ctx, id, userID)
		if checkErr != nil {
			return false, checkErr
		}
		if exists {
			return false, fmt.Errorf("version_conflict")
		}
		return false, nil
	}
	return true, nil
}

// existsByIDAndUser checks if a todo exists for the given user.
func (r *Repository) existsByIDAndUser(ctx context.Context, id, userID int64) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM todos WHERE id = $1 AND user_id = $2)`,
		id, userID,
	).Scan(&exists)
	return exists, err
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
