package handlers

import (
	"1337b04rd/internal/domain"
	"context"
	"mime/multipart"
)

func GetUserFromContext(ctx context.Context) (*domain.User, bool) {
	user, ok := ctx.Value(userContextKey).(*domain.User)
	return user, ok
}

// func PutCommentsInsideComments(post *domain.Post, comments []*domain.Comment) {
	
// }

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
