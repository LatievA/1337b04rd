package handlers

import (
	"context"
	"log/slog"
	"net/http"
)

type contextKey string

const userContextKey contextKey = "user"

func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		sessionToken := ""
		cookie, err := r.Cookie("session_token")
		if err != nil {
			slog.Debug("No session cookie found, will create new session")
		} else {
			sessionToken = cookie.Value
		}

		user, isNew, err := h.userService.GetOrCreateUser(ctx, sessionToken)
		if err != nil {
			slog.Error("Failed to get user by session", "err", err)
			h.HandleHTTPError(w, r, "Internal server error", http.StatusInternalServerError)
			return
		}

		if isNew {
			http.SetCookie(w, &http.Cookie{
				Name:     "session_token",
				Value:    user.SessionToken,
				Path:     "/",
				HttpOnly: true,
				Secure:   true, // Add this for HTTPS
				SameSite: http.SameSiteStrictMode,
				MaxAge:   7 * 24 * 60 * 60, // 7 days
			})
		}

		// Add user to context for downstream handlers
		ctx = context.WithValue(ctx, userContextKey, user)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
