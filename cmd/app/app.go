package app

import (
	"1337b04rd/internal/adapters/db/repository"
	"1337b04rd/internal/adapters/external_api"
	"1337b04rd/internal/adapters/handlers"
	"1337b04rd/internal/config"
	"1337b04rd/internal/logger"
	"1337b04rd/internal/server"
	"1337b04rd/internal/services"
	"log/slog"
)

func RunServer() {
	logger.Init(slog.LevelDebug)

	config, err := config.NewConfig()
	if err != nil {
		slog.Error("Failed to get configures", "error", err)
		return
	}

	db, err := repository.ConnectToDB(config.DBConfig)
	if err != nil {
		slog.Error("Failed to connect to database.", "error", err)
		return
	}
	postRepo := repository.NewPostRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	userRepo := repository.NewUserRepository(db)

	avatarProvider := external_api.NewRickAndMortyClient()

	userService := services.NewUserService(userRepo, avatarProvider)
	postService := services.NewPostService(postRepo, commentRepo)
	commentService := services.NewCommentService(commentRepo, postRepo)

	handler := handlers.NewHandler(userService, postService, commentService)
	server := server.NewServer(config, handler)

	server.Run()
}
