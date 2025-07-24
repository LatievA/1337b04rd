package handlers

import (
	"1337b04rd/internal/domain"
	"encoding/json"
	"log/slog"
	"net/http"
)

type PostHandler struct {
	postService domain.PostService
}

func NewPostHandler(postService domain.PostService) *PostHandler {
	return &PostHandler{postService: postService}
}

func (h *PostHandler) PostRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /posts", h.ListPosts)
	mux.HandleFunc("GET /archive", h.ListArchivedPosts)
	mux.HandleFunc("POST /post", h.CreatePost)
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	op := "POST /post"
	ctx := r.Context()

	var req CreatePostRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Invalid request body", "op", op, "error", err)
		RespondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON format"})
		return
	}

	post, err := h.postService.CreatePost(ctx, ???, req.Title, req.Content, req.ImageURL)
	if err != nil {
		slog.Error("Failed to create post", "op", op, "error", err)
		RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create post"})
		return
	}

	slog.Info("Post created successfully", "op", op, "postID", post.ID)
	RespondJSON(w, http.StatusCreated, post)
}

func (h *PostHandler) ListPosts(w http.ResponseWriter, r *http.Request) {
	op := "GET /posts"
	ctx := r.Context()

	posts, err := h.postService.ListPosts(ctx, false)
	if err != nil {
		slog.Error("Failed to get posts!", "OP", op, "error", err)
		RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
		return
	}
	slog.Info("Posts extracted succesfully!", "OP", op)
	RespondJSON(w, http.StatusOK, posts)
}

func (h *PostHandler) ListArchivedPosts(w http.ResponseWriter, r *http.Request) {
	op := "GET /archive"
	ctx := r.Context()

	posts, err := h.postService.ListPosts(ctx, true)
	if err != nil {
		slog.Error("Failed to get posts!", "OP", op, "error", err)
		RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
		return
	}
	slog.Info("Posts extracted succesfully!", "OP", op)
	RespondJSON(w, http.StatusOK, posts)
}
