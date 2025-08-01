package handlers

import (
	"1337b04rd/internal/domain"
	"net/http"
)

type Handler struct {
	userService    domain.UserService
	postService    domain.PostService
	commentService domain.CommentService
	s3Service      domain.S3Service
}

func NewHandler(userService domain.UserService, postService domain.PostService, commentService domain.CommentService, s3Service domain.S3Service) *Handler {
	return &Handler{
		userService:    userService,
		postService:    postService,
		commentService: commentService,
		s3Service:      s3Service,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("GET /catalog", h.AuthMiddleware(http.HandlerFunc(h.ListPosts)))
	mux.Handle("GET /archive", h.AuthMiddleware(http.HandlerFunc(h.ListArchivedPosts)))
	mux.Handle("GET /post/{id}", h.AuthMiddleware(http.HandlerFunc(h.GetPost)))
	mux.Handle("GET /archive-post/{id}", h.AuthMiddleware(http.HandlerFunc(h.GetArchivePost)))
	mux.Handle("GET /create-post", h.AuthMiddleware(http.HandlerFunc(h.CreatePostForm)))
	mux.Handle("POST /create-post", h.AuthMiddleware(http.HandlerFunc(h.CreatePost)))
	mux.Handle("POST /post/{id}/comment", h.AuthMiddleware(http.HandlerFunc(h.CreateComment)))
	mux.Handle("GET /error", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.HandleHTTPError(w, r, "An expected error occurred.", http.StatusInternalServerError)
	}))
}
