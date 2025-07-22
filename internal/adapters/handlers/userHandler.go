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
	ctx := r.Context()

	cookie, err := r.Cookie("session_token")
	if err != nil {
		slog.Warn("Failed to get cookies", "err", err)
	}
	sessionToken := cookie.Value
	user, isNew, err := h.userService.GetOrCreateUser(ctx, sessionToken)
	if err != nil {
		slog.Error("Failed to get user by session", "err", err)
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
}
