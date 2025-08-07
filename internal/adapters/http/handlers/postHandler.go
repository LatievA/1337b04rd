package handlers

import (
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"path"
	"strconv"
	"time"
)

type PostFormData struct {
	Name     string
	Title    string
	Content  string
	ImageURL string
}

type TemplateData struct {
	FormData PostFormData
	Error    map[string]string
}

func (h *Handler) ListPosts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	posts, err := h.postService.ListPosts(ctx, false)
	if err != nil {
		slog.Error("Failed to fetch posts", "err", err)
		h.HandleHTTPError(w, r, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("internal/ui/templates/catalog.html")
	if err != nil {
		slog.Error("Failed to parse template", "err", err)
		h.HandleHTTPError(w, r, "Could not load page", http.StatusInternalServerError)
		return
	}
	for _, post := range posts {
		post.User, err = h.userService.GetUserByID(ctx, post.UserID)
		if err != nil {
			slog.Error("Failed to fetch post user", "err", err)
			h.HandleHTTPError(w, r, "Failed to fetch post user", http.StatusInternalServerError)
			return
		}
	}

	err = tmpl.Execute(w, posts)
	if err != nil {
		slog.Error("Failed to execute template", "err", err)
		h.HandleHTTPError(w, r, "Could not load page", http.StatusInternalServerError)
		return
	}
} // Works correctly

func (h *Handler) ListArchivedPosts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	posts, err := h.postService.ListPosts(ctx, true)
	if err != nil {
		slog.Error("Failed to fetch posts", "err", err)
		h.HandleHTTPError(w, r, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("internal/ui/templates/archive.html")
	if err != nil {
		slog.Error("Failed to parse template", "err", err)
		h.HandleHTTPError(w, r, "Could not load page", http.StatusInternalServerError)
		return
	}

	for _, post := range posts {
		post.User, err = h.userService.GetUserByID(ctx, post.UserID)
		if err != nil {
			slog.Error("Failed to fetch post user", "err", err)
			h.HandleHTTPError(w, r, "Failed to fetch post user", http.StatusInternalServerError)
			return
		}
	}
	err = tmpl.Execute(w, posts)
	if err != nil {
		slog.Error("Failed to execute template", "err", err)
		h.HandleHTTPError(w, r, "Could not load page", http.StatusInternalServerError)
		return
	}
} // Works correctly

func (h *Handler) GetPost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	postID := path.Base(r.URL.Path)
	if postID == "" {
		h.HandleHTTPError(w, r, "Post ID is required", http.StatusBadRequest)
		return
	}
	ipostID, err := strconv.Atoi(postID)
	if err != nil {
		slog.Error("Invalid post ID", "err", err)
	}

	post, err := h.postService.GetPostByID(ctx, ipostID)
	if err != nil {
		slog.Error("Failed to fetch post", "err", err)
		h.HandleHTTPError(w, r, "Failed to fetch post", http.StatusInternalServerError)
		return
	}
	post.User, err = h.userService.GetUserByID(ctx, post.UserID)
	if err != nil {
		slog.Error("Failed to fetch post user", "err", err)
		h.HandleHTTPError(w, r, "Failed to fetch post user", http.StatusInternalServerError)
		return
	}

	for _, comment := range post.Comments {
		comment.User, _ = h.userService.GetUserByID(ctx, comment.UserID)
		if comment.User == nil {
			slog.Error("Failed to fetch comment user", "commentID", comment.ID)
			h.HandleHTTPError(w, r, "Failed to fetch comment user", http.StatusInternalServerError)
			return
		}
		if comment.ParentID > 0 {
			for _, comment1 := range post.Comments {
				if comment1.ID == comment.ParentID {
					comment1.Comments = append(comment1.Comments, comment)
					break
				}
			}
		}
	}

	tmpl, err := template.ParseFiles("internal/ui/templates/post.html")
	if err != nil {
		slog.Error("Failed to parse template", "err", err)
		h.HandleHTTPError(w, r, "Could not load page", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, post)
	if err != nil {
		slog.Error("Failed to execute template", "err", err)
		h.HandleHTTPError(w, r, "Could not load page", http.StatusInternalServerError)
		return
	}
} // Works correctly

func (h *Handler) GetArchivePost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	postID := path.Base(r.URL.Path)
	if postID == "" {
		h.HandleHTTPError(w, r, "Post ID is required", http.StatusBadRequest)
		return
	}
	ipostID, err := strconv.Atoi(postID)
	if err != nil {
		slog.Error("Invalid post ID", "err", err)
	}

	post, err := h.postService.GetPostByID(ctx, ipostID)
	if err != nil {
		slog.Error("Failed to fetch post", "err", err)
		h.HandleHTTPError(w, r, "Failed to fetch post", http.StatusInternalServerError)
		return
	}
	post.User, err = h.userService.GetUserByID(ctx, post.UserID)
	if err != nil {
		slog.Error("Failed to fetch post user", "err", err)
		h.HandleHTTPError(w, r, "Failed to fetch post user", http.StatusInternalServerError)
		return
	}

	for _, comment := range post.Comments {
		comment.User, _ = h.userService.GetUserByID(ctx, comment.UserID)
		if comment.User == nil {
			slog.Error("Failed to fetch comment user", "commentID", comment.ID)
			h.HandleHTTPError(w, r, "Failed to fetch comment user", http.StatusInternalServerError)
			return
		}
		if comment.ParentID > 0 {
			for _, comment1 := range post.Comments {
				if comment1.ID == comment.ParentID {
					comment1.Comments = append(comment1.Comments, comment)
					break
				}
			}
		}
	}

	tmpl, err := template.ParseFiles("internal/ui/templates/archive-post.html")
	if err != nil {
		slog.Error("Failed to parse template", "err", err)
		h.HandleHTTPError(w, r, "Could not load page", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, post)
	if err != nil {
		slog.Error("Failed to execute template", "err", err)
		h.HandleHTTPError(w, r, "Could not load page", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) CreatePostForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	_, ok := GetUserFromContext(ctx)
	if !ok {
		h.HandleHTTPError(w, r, "Unauthorized", http.StatusUnauthorized)
		return
	}

	tmpl, err := template.ParseFiles("internal/ui/templates/create-post.html")
	if err != nil {
		slog.Error("Failed to parse template", "err", err)
		h.HandleHTTPError(w, r, "Could not load page", http.StatusInternalServerError)
		return
	}

	data := TemplateData{
		FormData: PostFormData{},
		Error:    make(map[string]string),
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		slog.Error("Failed to execute template", "err", err)
		h.HandleHTTPError(w, r, "Could not load page", http.StatusInternalServerError)
		return
	}
} // Works correctly

func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := GetUserFromContext(ctx)
	if !ok {
		h.HandleHTTPError(w, r, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err := r.ParseMultipartForm(10 << 20) // 10MB max memory
	if err != nil {
		h.HandleHTTPError(w, r, "Unable to parse form", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	title := r.FormValue("title")
	content := r.FormValue("content")

	// Validation
	errors := make(map[string]string)
	if title == "" {
		errors["title"] = "Title is required"
	}
	if content == "" {
		errors["content"] = "Content is required"
	}

	// If validation fails, re-display form with errors
	if len(errors) > 0 {
		data := TemplateData{
			FormData: PostFormData{
				Name:    name,
				Title:   title,
				Content: content,
			},
			Error: errors,
		}

		tmpl := template.Must(template.ParseFiles("create-post.html"))
		tmpl.Execute(w, data)
		return
	}

	// Test
	if name != "" {
		h.userService.UpdateUserName(ctx, user.ID, name)
		user.Name = name // Update user in context after changing name
	}

	// Add triple-s implemenatation for file upload
	var imageURL string
	file, fh, err := r.FormFile("image")
	if err == nil && fh != nil {
		defer file.Close()

		raw, err := io.ReadAll(file)
		if err != nil {
			slog.Error("Failed to read image", "err", err)
			h.HandleHTTPError(w, r, "Failed to read image", http.StatusInternalServerError)
			return
		}

		key := fmt.Sprintf("%d-%s", time.Now().Unix(), fh.Filename)

		url, err := h.s3Service.UploadImage(ctx, raw, "posts", key)
		if err != nil {
			slog.Error("Failed to upload image to S3", "err", err)
			h.HandleHTTPError(w, r, "Failed to upload to S3", http.StatusInternalServerError)
			return
		}
		imageURL = url
	}
	_, err = h.postService.CreatePost(ctx, user.ID, user.Name, title, content, imageURL)
	if err != nil {
		slog.Error("Failed to create post", "err", err)
		h.HandleHTTPError(w, r, "Failed to create post", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/catalog", http.StatusSeeOther)
} // Works correctly
