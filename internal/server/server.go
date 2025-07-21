package server

import (
	"1337b04rd/internal/adapters/handlers"
	"1337b04rd/internal/config"
)

type Server struct {
	httpServer *http.Server 
}

func NewServer(config *config.Config, handler *handlers.Handler) *Server {
	mux := http.NewServeMux()
	handler.User.RegisterRoutes(mux)
	

	srv := http.Server{
		Addr:
	}
	
	return &Server{
		httpServer: srv,
	}
}