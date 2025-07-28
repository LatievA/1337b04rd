package handlers

import (
	"1337b04rd/internal/domain"
	"context"
	"log/slog"
	"time"
	"mime/multipart"
)

func GetUserFromContext(ctx context.Context) (*domain.User, bool) {
	user, ok := ctx.Value(userContextKey).(*domain.User)
	return user, ok
}

func (h *Handler) StartArchiveWorker() {
	go func() {
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()

		slog.Info("Archive worker started", "interval", "15s")
		
		for range ticker.C {
			slog.Info("Archiving expired posts")
			err:= h.postService.ArchiveOldPosts(context.Background())
			if err != nil {
				slog.Error("Failed to archive expired posts", "err", err)
				continue
			} else {
				slog.Info("Expired posts archived successfully")
			}
		}
	}()
}

func isValidImageType(fileHeader *multipart.FileHeader) bool {
	allowedTypes := []string{"image/jpeg", "image/png", "image/gif", "image/webp"}
	contentType := fileHeader.Header.Get("Content-Type")

	for _, allowedType := range allowedTypes {
		if contentType == allowedType {
			return true
		}
	}
	return false
}
