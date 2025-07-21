package handlers

import (
	"1337b04rd/internal/domain"
)

type PostHandler struct{
	postService domain.PostService
}

func NewPostHandler(postService domain.PostService) *PostHandler {
	return &PostHandler{postService: postService}
}
