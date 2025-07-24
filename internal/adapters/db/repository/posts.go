package repository

import (
	"context"
	"database/sql"
	"fmt"

	"1337b04rd/internal/domain"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) domain.PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) Save(ctx context.Context, post *domain.Post) (int, error) {
	var postID int
	query := `INSERT INTO posts(session_id, title, content, image_url)
			  VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, post.UserID, post.Title, post.Content, post.ImageURL).Scan(&postID)
	if err != nil {
		return -1, err
	}

	return postID, nil
}

func (r *PostRepository) FindByID(ctx context.Context, id int) (*domain.Post, error) {
	post := &domain.Post{}

	query := `SELECT * FROM posts WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID, &post.UserID, &post.Title, &post.Content, &post.ImageURL, &post.CreatedAt, &post.ArchivedAt, &post.Archived)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("post with id %d doesn't exist", id)
		}
		return nil, err
	}

	queryComments := `SELECT * FROM comments WHERE post_id = $1`
	rows, err := r.db.QueryContext(ctx, queryComments, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return post, nil
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		comment := &domain.Comment{}
		err = rows.Scan(&comment.ID, &comment.UserID, &comment.PostID, &comment.ParentID, &comment.Content, &comment.CreatedAt)
		if err != nil {
			return nil, err
		}
		post.Comments = append(post.Comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return post, err
}

func (r *PostRepository) FindAll(ctx context.Context, archived bool) ([]*domain.Post, error) {
	posts := []*domain.Post{}
	query := `SELECT * FROM posts WHERE archived = $1`

	rows, err := r.db.QueryContext(ctx, query, archived)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		post := &domain.Post{}
		err = rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.ImageURL, &post.CreatedAt, &post.ArchivedAt, &post.Archived)
		if err != nil {
			return nil, err
		}

		queryComments := `SELECT * FROM comments WHERE post_id = $1`
		rowsCom, err := r.db.QueryContext(ctx, queryComments, post.ID)
		if err != nil {
			if err == sql.ErrNoRows {
				continue
			}
			return nil, err
		}
		defer rowsCom.Close()

		for rowsCom.Next() {
			comment := &domain.Comment{}
			err = rowsCom.Scan(&comment.ID, &comment.UserID, &comment.PostID, &comment.ParentID, &comment.Content, &comment.CreatedAt)
			if err != nil {
				return nil, err
			}
			post.Comments = append(post.Comments, comment)
		}
		if err := rowsCom.Err(); err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *PostRepository) Update(ctx context.Context, post *domain.Post) error {
	query := `UPDATE posts SET title = $1, content = $2, image_url = $3, archived_at = $4, is_archived = $5 WHERE id = $6`
	result, err := r.db.ExecContext(ctx, query, post.Title, post.Content, post.ImageURL, post.ArchivedAt, post.Archived, post.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no post found with id %d", post.ID)
	}
	return nil
}

func (r *PostRepository) ArchiveExpired(ctx context.Context) error {
	return nil
}
