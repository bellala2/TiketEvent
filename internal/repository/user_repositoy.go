package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/bellala2/TiketEvent/internal/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create menyimpan user baru. id, created_at, updated_at diisi oleh database.
func (r *UserRepository) Create(ctx context.Context, u *model.User) error {
	const query = `
		INSERT INTO users (name, email, password_hash, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query, u.Name, u.Email, u.PasswordHash, u.Role).
		Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		// 23505 = unique_violation di PostgreSQL (email sudah ada).
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrDuplicateEmail
		}
		return fmt.Errorf("insert user: %w", err)
	}
	return nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	const query = `
		SELECT id, name, email, password_hash, role, created_at, updated_at
		FROM users
		WHERE LOWER(email) = LOWER($1)`
	return r.scanOne(r.db.QueryRowContext(ctx, query, email))
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	const query = `
		SELECT id, name, email, password_hash, role, created_at, updated_at
		FROM users
		WHERE id = $1`
	return r.scanOne(r.db.QueryRowContext(ctx, query, id))
}

func (r *UserRepository) scanOne(row *sql.Row) (*model.User, error) {
	var u model.User
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("scan user: %w", err)
	}
	return &u, nil
}
