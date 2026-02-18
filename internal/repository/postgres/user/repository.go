package user

import (
	user2 "ProjectManagementAPI/internal/domain/user"
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

type Repository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, u *user2.User) error {
	const query = `INSERT INTO users(id, email, name) VALUES ($1, $2, $3)`

	_, err := r.db.ExecContext(ctx, query, u.ID, u.Email, u.Name)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return user2.ErrEmailAlreadyExists
		}
		return err
	}

	return nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*user2.User, error) {
	const query = `SELECT id, email, name FROM users WHERE id=$1`

	u := &user2.User{}
	err := r.db.QueryRowContext(ctx, query, id).
		Scan(&u.ID, &u.Email, &u.Name)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, user2.ErrUserNotFound
	}

	return u, err
}

func (r *Repository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	const query = `DELETE FROM users WHERE id=$1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
