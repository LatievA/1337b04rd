package server

import (
	"1337b04rd/internal/adapters/handlers"
	"1337b04rd/internal/config"
	"fmt"
	"log/slog"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(config *config.Config, handler *handlers.Handler) *Server {
	mux := http.NewServeMux()
	handler.User.UserRoutes(mux)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":s", config.ServerConfig.Port),
		Handler: mux,
	}

	return &Server{
		httpServer: srv,
	}
}

func (s *Server) Run() error {
	slog.Info("Server is starting", "port", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}
