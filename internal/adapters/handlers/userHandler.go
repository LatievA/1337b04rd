package handlers

import (
	"1337b04rd/internal/domain"
	"log/slog"
	"net/http"
)

type UserHandler struct {
	userService domain.UserService
}

func NewUserHandler(userService domain.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) UserRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /", h.HandleSession)
}

func (h *UserHandler) HandleSession(w http.ResponseWriter, r *http.Request) {
	// Redirect to /catalog
	http.Redirect(w, r, "/catalog", http.StatusSeeOther)

	ctx := r.Context()

	cookie, err := r.Cookie("session_token")
	var sessionToken string
	if err != nil {
		slog.Warn("Failed to get cookies", "err", err)
	}
	if cookie != nil {
		sessionToken = cookie.Value
	}

	err = nil
	user, isNew, err := h.userService.GetOrCreateUser(ctx, sessionToken)
	if err != nil {
		slog.Error("Failed to get user by session", "err", err)
		return
	}

	if isNew {
		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    user.SessionToken,
			Path:     "/",
			HttpOnly: true,
			MaxAge:   7 * 24 * 60 * 60,
		})
	}
}
