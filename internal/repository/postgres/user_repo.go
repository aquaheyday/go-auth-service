package postgres

import (
	"context"
	"database/sql"
	"github.com/aquaheyday/go-auth-service/internal/domain"
	"github.com/google/uuid"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) (string, error) {
	id := uuid.New().String()
	query := `INSERT INTO users (id, email, password_hash, created_at) VALUES ($1, $2, $3, NOW())`
	if _, err := r.db.ExecContext(ctx, query, id, user.Email, user.PasswordHash); err != nil {
		return "", err
	}
	return id, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var u domain.User
	query := `SELECT id, email, password_hash, created_at FROM users WHERE email = $1`
	row := r.db.QueryRowContext(ctx, query, email)
	if err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}
