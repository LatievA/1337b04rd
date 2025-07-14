package repository

import (
	"1337b04rd/internal/domain"
	"context"
	"database/sql"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) Save(ctx context.Context, post *domain.Post) (int, error) {
	var postID int
	query := `INSERT INTO posts(session_id, title, content, image_url)
			  VALUES ($1, $2, $3) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, post.UserID, post.Title, post.Content, post.ImageURL).Scan(&postID)
	if err != nil {
		return -1, err
	}
	
	return postID, nil
}

func (r *PostRepository) FindByID(ctx context.Context, id int) (*domain.Post, error) {
	post := &domain.Post{}
	query := `SELECT * FROM posts WHERE id = $1`
	
	// need to add COMMENTS recieving
	
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID, &post.UserID, &post.Title, &post.Content, &post.ImageURL, &post.CreatedAt, &post.ArchivedAt, &post.Archived)
	if err != nil {
		return nil, err
	}
	return post, err
}

func (r *PostRepository) FindAll(ctx context.Context, archived bool) ([]*domain.Post, error) {
	posts := []*domain.Post{}
	query := `SELECT * FROM posts`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	// need to add COMMENTS recieving
	for rows.Next() {
		post := &domain.Post{}
		rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.ImageURL, &post.CreatedAt, &post.ArchivedAt, &post.Archived)
		posts = append(posts, post)
	}
	
}