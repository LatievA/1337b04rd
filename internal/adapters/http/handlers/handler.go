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
		userService: userService,
		postService: postService,
		commentService: commentService,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /", h.HandleSession)
	mux.HandleFunc("PUT /user", h.UpdateUserName)
	mux.HandleFunc("GET /posts", h.ListPosts)
	mux.HandleFunc("GET /archive", h.ListArchivedPosts)
	mux.HandleFunc("POST /post", h.CreatePost)
}


func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	op := "POST /post"
	ctx := r.Context()

	var req CreatePostRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Invalid request body", "op", op, "error", err)
		RespondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON format"})
		return
	}

	// How to implement that
	post, err := h.postService.CreatePost(ctx, ???, req.Title, req.Content, req.ImageURL)
	if err != nil {
		slog.Error("Failed to create post", "op", op, "error", err)
		RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create post"})
		return
	}

	slog.Info("Post created successfully", "op", op, "postID", post.ID)
	RespondJSON(w, http.StatusCreated, post)
}

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


func (h *Handler) HandleSession(w http.ResponseWriter, r *http.Request) {
	op := "GET /"
	ctx := r.Context()

	cookie, err := r.Cookie("session_token")
	var sessionToken string
	if err != nil {
		slog.Warn("Failed to get cookies", "OP", op, "err", err)
	}
	if cookie != nil {
		sessionToken = cookie.Value
	}

	err = nil
	user, isNew, err := h.userService.GetOrCreateUser(ctx, sessionToken)
	if err != nil {
		slog.Error("Failed to get user by session", "OP", op, "err", err)
		return
	}

	if isNew {
		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    user.Session,
			Path:     "/",
			HttpOnly: true,
			MaxAge:   7 * 24 * 60 * 60,
		})
	}

	tmpl, err := template.ParseFiles("internal/ui/templates/catalog.html")
	if err != nil {
		slog.Error("Failed to parse template", "OP", op, "err", err)
		http.Error(w, "Could not load page", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, "")
	if err != nil {
		slog.Error("Failed to execute template", "OP", op, "err", err)
		http.Error(w, "Could not load page", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UpdateUserName(w http.ResponseWriter, r *http.Request) {
}
