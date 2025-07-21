package handlers

import (
	"1337b04rd/internal/domain"
)

type Handler struct {
	User    *UserHandler
	Post    *PostHandler
	Comment *CommentHandler
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
