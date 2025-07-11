package services

import (
	"1337b04rd/internal/domain"
	"context"
	"errors"
	"time"
)

type postService struct {
	postRepo domain.PostRepository
	commentRepo domain.CommentRepository
}

type commentService struct {
	commentRepo domain.CommentRepository
	postRepo    domain.PostRepository
}

type userService struct {
	userRepo         domain.UserRepository
	rickAndMortyAPI  domain.RickAndMortyAPI
}

func NewPostService(postRepo domain.PostRepository, commentRepo domain.CommentRepository) domain.PostService {
	return &postService{postRepo: postRepo, commentRepo: commentRepo}
}

func NewCommentService(commentRepo domain.CommentRepository, postRepo domain.PostRepository) domain.CommentService {
	return &commentService{commentRepo: commentRepo, postRepo: postRepo}
}

func NewUserService(userRepo domain.UserRepository, api domain.RickAndMortyAPI) domain.UserService {
	return &userService{userRepo: userRepo, rickAndMortyAPI: api}
}

func (s *postService) CreatePost(ctx context.Context, userID int, title, content string, imageURL *string) (*domain.Post, error) {
	post := &domain.Post{
		UserID:    userID,
		Title:     title,
		Content:   content,
		ImageURL:  imageURL,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	id, err := s.postRepo.Save(ctx, post)
	if err != nil {
		return nil, err
	}
	post.ID = id
	return post, nil
}

func (s *postService) GetPostByID(ctx context.Context, postID int) (*domain.Post, error) {
	return s.postRepo.FindByID(ctx, postID)
}

func (s *postService) ListPosts(ctx context.Context, archived bool) ([]*domain.Post, error) {
	return s.postRepo.FindAll(ctx, archived)
}

func (s *postService) ArchiveOldPosts(ctx context.Context) error {
	return s.postRepo.ArchiveExpired(ctx)
}

func (s *commentService) AddComment(ctx context.Context, userID, postID, parentID int, content string) (*domain.Comment, error) {
	_, err := s.postRepo.FindByID(ctx, postID)
	if err != nil {
		return nil, errors.New("post not found")
	}
	comment := &domain.Comment{
		UserID:    userID,
		ParentID:  parentID,
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

func (s *commentService) GetCommentsByPostID(ctx context.Context, postID int) ([]*domain.Comment, error) {
	return s.commentRepo.FindByPostID(ctx, postID)
}

func (s *userService) GetOrCreateUser(ctx context.Context, sessionID string) (*domain.User, error) {
	user, err := s.userRepo.FindBySessionID(ctx, sessionID)
	if err == nil {
		return user, nil
	}

	name, avatarURL, err := s.rickAndMortyAPI.GetRandomCharacter(ctx)
	if err != nil {
		return nil, err
	}

	newUser := &domain.User{
		Name:      name,
		AvatarURL: &avatarURL,
	}
	id, err := s.userRepo.Save(ctx, newUser)
	if err != nil {
		return nil, err
	}
	newUser.ID = id
	return newUser, nil
}

func (s *userService) UpdateUserName(ctx context.Context, userID int, newName string) error {
	return s.userRepo.UpdateName(ctx, userID, newName)
}
