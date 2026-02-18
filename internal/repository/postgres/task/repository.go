package task

import (
	task2 "ProjectManagementAPI/internal/domain/task"
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Repository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, t *task2.Task) error {
	t.ID = uuid.New()
	t.CreatedAt = time.Now()

	const query = `INSERT INTO tasks(id, title, description, status, created_at) VALUES($1,$2,$3,$4,$5)`
	_, err := r.db.ExecContext(ctx, query, t.ID, t.Title, t.Description, t.Status, t.CreatedAt)

	for _, userID := range t.Assignees {
		const linkQuery = `INSERT INTO user_tasks(user_id, task_id) VALUES($1, $2)`
		if _, err := r.db.ExecContext(ctx, linkQuery, userID, t.ID); err != nil {
			return err
		}
	}

	return err
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*task2.Task, error) {
	const taskQuery = `SELECT id, title, description, status, created_at FROM tasks WHERE id=$1`

	t := &task2.Task{}
	err := r.db.QueryRowContext(ctx, taskQuery, id).Scan(
		&t.ID, &t.Title, &t.Description, &t.Status, &t.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, task2.ErrTaskNotFound
	} else if err != nil {
		return nil, err
	}

	// Подгружаем исполнителей
	const assigneeQuery = `SELECT user_id FROM user_tasks WHERE task_id=$1`
	rows, err := r.db.QueryContext(ctx, assigneeQuery, id)
	if err != nil {
		return nil, err
	}

	var assignees []uuid.UUID
	for rows.Next() {
		var userID uuid.UUID
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		assignees = append(assignees, userID)
	}
	t.Assignees = assignees

	return t, nil
}

func (r *Repository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	const query = `DELETE FROM tasks WHERE id=$1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
