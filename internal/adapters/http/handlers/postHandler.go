package handlers

import (
	"log/slog"
	"net/http"
)

func (h *Handler) ListPosts(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) ListArchivedPosts(w http.ResponseWriter, r *http.Request) {
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
