package services

import (
	"1337b04rd/internal/domain"
	"context"
	"fmt"
	"time"
)

type PostService struct {
	postRepo    domain.PostRepository
	commentRepo domain.CommentRepository
	userRepo    domain.UserRepository
	s3Service   domain.S3Service
}

func NewPostService(postRepo domain.PostRepository, commentRepo domain.CommentRepository, userRepo domain.UserRepository) domain.PostService {
	return &PostService{postRepo: postRepo, commentRepo: commentRepo, userRepo: userRepo}
}

func (s *PostService) CreatePost(ctx context.Context, userID int, name, title, content, imageURL string) (*domain.Post, error) {

	post := &domain.Post{
		UserID:     userID,
		Username:   name,
		Title:      title,
		Content:    content,
		ImageURL:   imageURL,
		CreatedAt:  time.Now(),
		ArchivedAt: time.Now().Add(15 * time.Minute),
	}
	id, err := s.postRepo.Save(ctx, post)
	if err != nil {
		return nil, fmt.Errorf("failed to save created post: %s", err)
	}
	post.ID = id
	return post, nil
}

func (s *PostService) GetPostByID(ctx context.Context, postID int) (*domain.Post, error) {
	return s.postRepo.FindByID(ctx, postID)
}

func (s *PostService) ListPosts(ctx context.Context, archived bool) ([]*domain.Post, error) {
	return s.postRepo.FindAll(ctx, archived)
}

func (s *PostService) ArchiveOldPosts(ctx context.Context) error {
	return s.postRepo.ArchiveExpired(ctx)
}