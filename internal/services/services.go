package services

import (
	"1337b04rd/internal/domain"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"time"
)

type PostService struct {
	postRepo    domain.PostRepository
	commentRepo domain.CommentRepository
}

type CommentService struct {
	commentRepo domain.CommentRepository
	postRepo    domain.PostRepository
}

type UserService struct {
	userRepo        domain.UserRepository
	rickAndMortyAPI domain.RickAndMortyAPI
}

func NewPostService(postRepo domain.PostRepository, commentRepo domain.CommentRepository) domain.PostService {
	return &PostService{postRepo: postRepo, commentRepo: commentRepo}
}

func NewCommentService(commentRepo domain.CommentRepository, postRepo domain.PostRepository) domain.CommentService {
	return &CommentService{commentRepo: commentRepo, postRepo: postRepo}
}

func NewUserService(userRepo domain.UserRepository, api domain.RickAndMortyAPI) domain.UserService {
	return &UserService{userRepo: userRepo, rickAndMortyAPI: api}
}

func (s *PostService) CreatePost(ctx context.Context, sessionToken, title, content string, imageURL *string) (*domain.Post, error) {
	

	post := &domain.Post{
		Username:     sessionToken,
		Title:      title,
		Content:    content,
		ImageURL:   imageURL,
		CreatedAt:  time.Now(),
		ArchivedAt: time.Now().Add(15 * time.Minute),
	}
	id, err := s.postRepo.Save(ctx, post)
	if err != nil {
		return nil, err
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

func (s *UserService) GetOrCreateUser(ctx context.Context, sessionToken string) (*domain.User, bool, error) {
	var isNew bool
	slog.Info(sessionToken)
	if sessionToken == "" {
		newSessionToken, err := s.generateSession()
		isNew = true
		if err != nil {
			return nil, isNew, fmt.Errorf("failed to create new session token")
		}
		sessionToken = newSessionToken
	} else {
		user, err := s.userRepo.FindBySessionToken(ctx, sessionToken)
		if err == nil {
			slog.Info("Returned user from db")
			return user, isNew, nil
		}
		isNew = true
	}

	name, avatarURL, err := s.rickAndMortyAPI.GetRandomCharacter(ctx)
	if err != nil {
		return nil, isNew, err
	}

	newUser := &domain.User{
		Name:      name,
		AvatarURL: &avatarURL,
		SessionToken:   sessionToken,
	}
	id, err := s.userRepo.Save(ctx, newUser)
	if err != nil {
		return nil, isNew, err
	}
	newUser.ID = id
	return newUser, isNew, nil
}

func (s *UserService) UpdateUserName(ctx context.Context, userID int, newName string) error {
	return s.userRepo.UpdateName(ctx, userID, newName)
}

func (s *UserService) generateSession() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}
