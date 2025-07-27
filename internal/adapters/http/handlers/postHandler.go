package handlers

import (
	"html/template"
	"log/slog"
	"net/http"
)

func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}
	sessionToken, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")
	imageURL := r.FormValue("image_url")

	var imageURLPtr *string
	if imageURL != "" {
		imageURLPtr = &imageURL
	}

	_, err = h.postService.CreatePost(ctx, sessionToken.Value, title, content, imageURLPtr)
	if err != nil {
		slog.Error("Failed to create post", "err", err)
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/catalog", http.StatusSeeOther)
}

func (h *Handler) ListPosts(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) ListArchivedPosts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	posts, err := h.postService.ListPosts(ctx, false)
	if err != nil {
		slog.Error("Failed to fetch posts", "err", err)
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("internal/ui/templates/archive.html")
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
