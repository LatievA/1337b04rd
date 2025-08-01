package handlers

import (
	"html/template"
	"log/slog"
	"net/http"
)

// Replace this function in your handlers
func (h *Handler)HandleHTTPError(w http.ResponseWriter, r *http.Request, message string, statusCode int) {
	slog.Error("HTTP Error", "status", statusCode, "message", message, "path", r.URL.Path)
	ServeErrorPage(w, r, statusCode, message)
}

func ServeErrorPage(w http.ResponseWriter, r *http.Request, statusCode int, message string) {
	w.WriteHeader(statusCode)

	tmpl, err := template.ParseFiles("internal/ui/templates/error.html")
	if err != nil {
		// Fallback if error template fails
		slog.Error("Failed to load error template", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := struct {
		StatusCode int
		Message    string
	}{
		StatusCode: statusCode,
		Message:    message,
	}

	if err := tmpl.Execute(w, data); err != nil {
		slog.Error("Failed to execute error template", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
