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

	var sessionID string
	cookie, err := r.Cookie("session_id")
	if err != nil {
		slog.Warn("Failed to get cookies", "err", err)
	}
	sessionToken := cookie.Value
	user, err := h.userService.GetOrCreateUser(ctx, sessionToken)
	if err != nil {
		slog.Error("Failed to get user by session", "err", err)
		return
	}
}
