package handlers

import (
	"log/slog"
	"net/http"
	"strconv"
)

func (h *Handler) SubmitComment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := GetUserFromContext(ctx)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	postID := r.FormValue("post_id")
	content := r.FormValue("content")

	if postID == "" || content == "" {
		http.Error(w, "Post ID and content are required", http.StatusBadRequest)
		return
	}
	ipostID, _ := strconv.Atoi(postID)
	iparentID, _ := strconv.Atoi(r.FormValue("parent_id"))

	_, err := h.commentService.AddComment(ctx, user.ID, ipostID, iparentID, content)
	if err != nil {
		slog.Error("Failed to submit comment", "err", err)
		http.Error(w, "Failed to submit comment", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/post/"+postID, http.StatusSeeOther)
}
