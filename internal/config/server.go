package config

import (
	"1337b04rd/internal/config"
	"net/http"
)

type Server struct {
	httpServer *http.Server 
}

func NewServer(config *config.Config) *Server {
	mux := http.NewServeMux()

	
	

	srv := http.Server{
		Addr:
	}
	
	return &Server{
		httpServer: srv,
	}
}