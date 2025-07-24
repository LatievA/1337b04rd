package handlers

import (
	"1337b04rd/internal/domain"
	"log/slog"
	"net/http"
	"text/template"
)

type UserHandler struct {
	userService domain.UserService
}

func NewUserHandler(userService domain.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) UserRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /", h.HandleSession)
	mux.HandleFunc("PUT /user", h.UpdateUserName)
}

func (h *UserHandler) HandleSession(w http.ResponseWriter, r *http.Request) {
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

func (h *UserHandler) UpdateUserName(w http.ResponseWriter, r *http.Request) {
}
