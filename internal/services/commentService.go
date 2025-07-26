package services

import (
	"1337b04rd/internal/domain"
	"context"
	"errors"
	"time"
)

type CommentService struct {
	commentRepo domain.CommentRepository
	postRepo    domain.PostRepository
}

func NewCommentService(commentRepo domain.CommentRepository, postRepo domain.PostRepository) domain.CommentService {
	return &CommentService{commentRepo: commentRepo, postRepo: postRepo}
}

func (s *CommentService) AddComment(ctx context.Context, userID, postID, parentID int, content string) (*domain.Comment, error) {
	_, err := s.postRepo.FindByID(ctx, postID)
	if err != nil {
		return nil, errors.New("post not found")
	}

	var parentPtr *int
	if parentID > 0 {
		parentPtr = &parentID
	}

	comment := &domain.Comment{
		UserID:    userID,
		ParentID:  parentPtr,
		Content:   content,
		CreatedAt: time.Now(),
	}
	id, err := s.commentRepo.Save(ctx, comment)
	if err != nil {
		return nil, err
	}
	comment.ID = id
	return comment, nil
}

func (s *CommentService) GetCommentsByPostID(ctx context.Context, postID int) ([]*domain.Comment, error) {
	return s.commentRepo.FindByPostID(ctx, postID)
}
