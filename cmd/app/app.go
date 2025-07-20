package app

import (
	"1337b04rd/internal/adapters/db/repository"
	"1337b04rd/internal/adapters/handlers"
	"1337b04rd/internal/config"
	"1337b04rd/internal/logger"
	"1337b04rd/internal/services"
	"fmt"
	"log"
	"log/slog"
	"net/http"
)

func RunServer() {
	logger.Init(slog.LevelDebug)
	db := 

	postRepo := repository.NewPostRepository()

	userService := services.NewUserService()
	postService := services.NewPostService()
	commentService := services.NewCommentService()

	handler := handlers.NewHandler(userService,postService,commentService)
	config, err := config.NewConfig()
	if err != nil {
		slog.Error("Failed to get configures", "error", err)
		return
	}
	
	addr := fmt.Sprintf(":%s", *port)
	log.Printf("Server is running on http://localhost%s\n", addr)

	if err := http.ListenAndServe(addr, handlers.RooterWays()); err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
}
