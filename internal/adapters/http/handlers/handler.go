package handlers

import (
	"1337b04rd/internal/domain"
)

type Handler struct {
	User    *UserHandler
	Post    *PostHandler
	Comment *CommentHandler
}

type CreatePostRequest struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	ImageURL string `json:"image_url"`
}

func NewHandler(
	userService domain.UserService,
	postService domain.PostService,
	commentService domain.CommentService,
) *Handler {
	return &Handler{
		User:    NewUserHandler(userService),
		Post:    NewPostHandler(postService),
		Comment: NewCommentHandler(commentService),
	}
}
