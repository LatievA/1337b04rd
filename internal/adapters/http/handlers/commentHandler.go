package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"path"
	"strconv"
)

func (h *Handler) CreateComment(w http.ResponseWriter, r *http.Request) {
	postID := path.Base(path.Dir(r.URL.Path))

	// Parse form data
	if err := r.ParseForm(); err != nil {
		slog.Error("Failed to parse form", "error", err)
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get form values
	name := r.FormValue("name")
	content := r.FormValue("content")
	replyTo := r.FormValue("reply_to")

	// Validate required fields
	if content == "" {
		http.Error(w, "Comment content is required", http.StatusBadRequest)
		return
	}

	// Get user info
	user, ok := GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Session required", http.StatusUnauthorized)
		return
	}

	if name != "" {
		user.Name = name
		h.userService.UpdateUserName(r.Context(), user.ID, name)
	}

	parentID := 0
	// Handle reply functionality
	if replyTo != "" {
		parentID, _ = strconv.Atoi(replyTo)
	}

	ipostID, err := strconv.Atoi(postID)
	if err != nil {
		slog.Error("Invalid post ID", "error", err)
	}

	// Save comment using repository
	_, err = h.commentService.AddComment(r.Context(), user.ID, ipostID, parentID, content)
	if err != nil {
		slog.Error("Failed to save comment", "error", err)
		http.Error(w, "Failed to save comment", http.StatusInternalServerError)
		return
	}

	// Add 15 minutes to post lifetime
	err = h.postService.AddTimeToPostLifetime(r.Context(), ipostID)

	// Redirect back to the post page
	http.Redirect(w, r, fmt.Sprintf("/post/%s", postID), http.StatusSeeOther)
} // Works correctly
