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
		db: db,
	}
}

func (r CommentRepository) Save(ctx context.Context, comment *domain.Comment) (int, error) {
    query := `
        INSERT INTO comments (session_id, post_id, parent_comment_id, content)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `

    var id int
    err := r.db.QueryRowContext(ctx, query,
        comment.UserID,
        comment.PostID,
        comment.ParentID,
        comment.Content,
    ).Scan(&id)

    if err != nil {
        return -1, err
    }

    return id, nil
}

func (r CommentRepository) FindByID(ctx context.Context, commentID int) (*domain.Comment, error) {
    var comment domain.Comment
    query := `
        Select *
        FROM comments
        WHERE id = $1
    `
    err := r.db.QueryRowContext(ctx, query, commentID).Scan(
        &comment.ID,
        &comment.UserID,
        &comment.PostID,
        &comment.ParentID,
        &comment.Content,
        &comment.CreatedAt,
    )
    if err != nil {
        return nil, err
    }
    return &comment, nil
}

func (r CommentRepository) FindByPostID(ctx context.Context, postID int) ([]*domain.Comment, error) {
    query := `
        SELECT id, session_id, post_id, parent_comment_id, content, created_at
        FROM comments
        WHERE post_id = $1
        ORDER BY created_at ASC
    `

    rows, err := r.db.QueryContext(ctx, query, postID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var comments []*domain.Comment
    for rows.Next() {
        var comment domain.Comment
        var parentCommentID sql.NullInt64
        var createdAt sql.NullTime
        
        err := rows.Scan(
            &comment.ID,
            &comment.UserID,
            &comment.PostID,
            &comment.ParentID,
            &comment.Content,
            &createdAt,
        )
        if err != nil {
            return nil, err
        }

        // Handle nullable fields
        if parentCommentID.Valid {
            comment.ParentID = int(parentCommentID.Int64)
        } else {
            comment.ParentID = 0 // or however you want to represent null
        }

        if createdAt.Valid {
            comment.CreatedAt = createdAt.Time
        }

        comments = append(comments, &comment)
    }

    if err = rows.Err(); err != nil {
        return nil, err
    }

    return comments, nil
}
