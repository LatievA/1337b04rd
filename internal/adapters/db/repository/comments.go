package repository

import (
	"1337b04rd/internal/domain"
	"context"
	"database/sql"
)

type CommentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) domain.CommentRepository {
	return &CommentRepository{
		db:db,
	}
}

func (r CommentRepository) Save(ctx context.Context, comment *domain.Comment) (int, error) {
	// need to implement
	return -1, nil
}

func (r CommentRepository) FindByPostID(ctx context.Context, postID int) ([]*domain.Comment, error) {
	// need to implement
	return nil, nil
}