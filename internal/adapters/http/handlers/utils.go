package handlers

import (
	"1337b04rd/internal/domain"
	"context"
)

func GetUserFromContext(ctx context.Context) (*domain.User, bool) {
	user, ok := ctx.Value("user").(*domain.User)
	return user, ok
}
