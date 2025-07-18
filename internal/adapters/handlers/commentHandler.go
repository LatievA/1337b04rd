package handlers

import "1337b04rd/internal/services"

type CommentHandler struct {
	commentService *services.CommentService
}