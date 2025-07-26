package domain

import (
	"context"
)

type PostService interface {
	CreatePost(ctx context.Context, userID int, title, content string, imageURL *string) (*Post, error)
	GetPostByID(ctx context.Context, postID int) (*Post, error)
	ListPosts(ctx context.Context, archived bool) ([]*Post, error)
	ArchiveOldPosts(ctx context.Context) error
}

type CommentService interface {
	AddComment(ctx context.Context, userID, postID, parentID int, content string) (*Comment, error)
	GetCommentsByPostID(ctx context.Context, postID int) ([]*Comment, error)
}

type UserService interface {
	GetOrCreateUser(ctx context.Context, sessionToken string) (*User, bool, error)
	UpdateUserName(ctx context.Context, userID int, newName string) error
}

type PostRepository interface {
	Save(ctx context.Context, post *Post) (int, error)
	FindByID(ctx context.Context, id int) (*Post, error)
	FindAll(ctx context.Context, archived bool) ([]*Post, error)
	Update(ctx context.Context, post *Post) error
	ArchiveExpired(ctx context.Context) error
}

type CommentRepository interface {
	Save(ctx context.Context, comment *Comment) (int, error)
	FindByPostID(ctx context.Context, postID int) ([]*Comment, error)
}

type UserRepository interface {
	FindBySessionToken(ctx context.Context, sessionToken string) (*User, error)
	Save(ctx context.Context, user *User) (int, error)
	UpdateName(ctx context.Context, userID int, newName string) error
}

type S3Service interface {
	UploadImage(ctx context.Context, fileData []byte, filename string) (string, error)
}

type RickAndMortyAPI interface {
	GetRandomCharacter(ctx context.Context) (name string, avatarURL string, err error)
}
