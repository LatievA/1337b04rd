package handlers

import (
	"1337b04rd/internal/domain"
	"log/slog"
	"net/http"
	"text/template"
)

type PostHandler struct {
	postService domain.PostService
}

func NewPostHandler(postService domain.PostService) *PostHandler {
	return &PostHandler{postService: postService}
}
func (h *PostHandler) PostRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /catalog", h.Catalog)
}

func (h *PostHandler) Catalog(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	posts, err := h.postService.ListPosts(ctx, false)
	if err != nil {
		slog.Error("Failed to fetch posts", "err", err)
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("internal/ui/templates/catalog.html")
	if err != nil {
		slog.Error("Failed to parse template", "err", err)
		http.Error(w, "Could not load page", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, posts)
	if err != nil {
		slog.Error("Failed to execute template", "err", err)
		http.Error(w, "Could not load page", http.StatusInternalServerError)
		return
	}
}
