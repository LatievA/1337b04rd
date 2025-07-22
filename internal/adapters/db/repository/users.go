package repository

import (
	"context"
	"database/sql"
	"fmt"

	"1337b04rd/internal/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindBySessionToken(ctx context.Context, sessionToken string) (*domain.User, error) {
	var user domain.User
	query := `SELECT * FROM user_sessions WHERE session_token = $1`
	err := r.db.QueryRowContext(ctx, query, sessionToken).Scan(
		&user.ID,
		&user.Session,
		&user.Name,
		&user.AvatarURL,
		&user.ExpiresAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Save(ctx context.Context, user *domain.User) (int, error) {
	var userID int
	query := `INSERT INTO user_sessions(session_token, name, avatar_url)
			  VALUES ($1, $2, $3) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, user.Session, user.Name, user.AvatarURL).Scan(&userID)
	if err != nil {
		return -1, err
	}

	return userID, nil
}

func (r *UserRepository) UpdateName(ctx context.Context, userID int, newName string) error {
	query := `UPDATE user_sessions SET name = $1 WHERE id = $2`

	result, err := r.db.ExecContext(ctx, query, newName, userID)
	if err != nil {
		return fmt.Errorf("failed to update name: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no user found with ID %d", userID)
	}

	return nil
}
