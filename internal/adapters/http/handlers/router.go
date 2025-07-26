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
	mux.Handle("GET /catalog.html", h.AuthMiddleware(http.HandlerFunc(h.ListPosts)))
	mux.Handle("GET /archive.html", h.AuthMiddleware(http.HandlerFunc(h.ListArchivedPosts)))
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