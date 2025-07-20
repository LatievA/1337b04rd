package handlers

import "1337b04rd/internal/domain"

type CommentHandler struct {
	commentService domain.CommentService
}

func NewCommentHandler(commentService domain.CommentService) *CommentHandler {
	return &CommentHandler{commentService: commentService}
}