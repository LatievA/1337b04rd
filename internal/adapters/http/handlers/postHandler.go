package handlers

import (
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"path"
)

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

func (h *Handler) GetPost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	postID := path.Base(r.URL.Path)
	if postID == "" {
		http.Error(w, "Post ID is required", http.StatusBadRequest)
		return
	}
	ipostID, err := strconv.Atoi(postID)
	if err != nil {
		slog.Error("Invalid post ID", "err", err)
	}

	post, err := h.postService.GetPostByID(ctx, ipostID)
	if err != nil {
		slog.Error("Failed to fetch post", "err", err)
		http.Error(w, "Failed to fetch post", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("internal/ui/templates/post.html")
	if err != nil {
		slog.Error("Failed to parse template", "err", err)
		http.Error(w, "Could not load page", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, post)
	if err != nil {
		slog.Error("Failed to execute template", "err", err)
		http.Error(w, "Could not load page", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetArchivePost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	postID := path.Base(r.URL.Path)
	if postID == "" {
		http.Error(w, "Post ID is required", http.StatusBadRequest)
		return
	}
	ipostID, err := strconv.Atoi(postID)
	if err != nil {
		slog.Error("Invalid post ID", "err", err)
	}

	post, err := h.postService.GetPostByID(ctx, ipostID)
	if err != nil {
		slog.Error("Failed to fetch post", "err", err)
		http.Error(w, "Failed to fetch post", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("internal/ui/templates/archive-post.html")
	if err != nil {
		slog.Error("Failed to parse template", "err", err)
		http.Error(w, "Could not load page", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, post)
	if err != nil {
		slog.Error("Failed to execute template", "err", err)
		http.Error(w, "Could not load page", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) CreatePostForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := GetUserFromContext(ctx)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	tmpl, err := template.ParseFiles("internal/ui/templates/create-post.html")
	if err != nil {
		slog.Error("Failed to parse template", "err", err)
		http.Error(w, "Could not load page", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, user)
	if err != nil {
		slog.Error("Failed to execute template", "err", err)
		http.Error(w, "Could not load page", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := GetUserFromContext(ctx)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err := r.ParseMultipartForm(32 << 20) // 32MB max memory
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	title := r.FormValue("title")
	content := r.FormValue("content")

	if name != "" {
		h.userService.UpdateUserName(ctx, user.ID, name)
		user.Name = name // Update user in context after changing name
	}

	if title == "" || content == "" {
		http.Error(w, "Title and content are required", http.StatusBadRequest)
		return
	}

	// Rewrite using triple-s
	var filePath string
	file, fileHeader, err := r.FormFile("file")
	if err == nil && fileHeader != nil {
		defer file.Close()

		// Create uploads directory if it doesn't exist
		os.MkdirAll("uploads", os.ModePerm)

		// Generate unique filename
		filename := fmt.Sprintf("%d_%s", time.Now().Unix(), fileHeader.Filename)
		filePath = filepath.Join("uploads", filename)

		// Save file
		outFile, err := os.Create(filePath)
		if err != nil {
			http.Error(w, "Unable to save file", http.StatusInternalServerError)
			return
		}
		defer outFile.Close()

		// Copy uploaded file to destination
		_, err = io.Copy(outFile, file)
		if err != nil {
			http.Error(w, "Unable to save file", http.StatusInternalServerError)
			return
		}
	}

	_, err = h.postService.CreatePost(ctx, user.ID, user.Name, title, content, nil)
	if err != nil {
		slog.Error("Failed to create post", "err", err)
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/catalog.html", http.StatusSeeOther)
}
