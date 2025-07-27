package handlers

import (
	"1337b04rd/internal/domain"
	"net/http"
)

type Handler struct {
	userService    domain.UserService
	postService    domain.PostService
	commentService domain.CommentService
}

func NewHandler(userService domain.UserService, postService domain.PostService, commentService domain.CommentService) *Handler {
	return &Handler{
		userService:    userService,
		postService:    postService,
		commentService: commentService,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("GET /catalog", h.AuthMiddleware(http.HandlerFunc(h.ListPosts)))
	mux.Handle("GET /archive", h.AuthMiddleware(http.HandlerFunc(h.ListArchivedPosts)))
	mux.Handle("GET /post/{id}", h.AuthMiddleware(http.HandlerFunc(h.GetPost)))
	mux.Handle("GET /archive-post/{id}", h.AuthMiddleware(http.HandlerFunc(h.GetArchivePost)))
	mux.Handle("GET /create-post", h.AuthMiddleware(http.HandlerFunc(h.CreatePostForm)))
	mux.Handle("POST /create-post", h.AuthMiddleware(http.HandlerFunc(h.CreatePost)))
	mux.Handle("POST /submit-comment", h.AuthMiddleware(http.HandlerFunc(h.SubmitComment)))
}

/*
user, ok := GetUserFromContext(r.Context())
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    
    // Use user in your logic
    slog.Info("User accessing posts", "userID", user.ID)
*/