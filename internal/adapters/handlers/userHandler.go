package handlers

import (
	"1337b04rd/internal/domain"
	"net/http"
)

type UserHandler struct {
	userService domain.UserService
}

func NewUserHandler(userService domain.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/user/me", h.HandleMe)
}

func (h *UserHandler) HandleMe(w http.ResponseWriter, r *http.Request) {
	// Пример: получить sessionToken из куки, вызвать h.userService.GetOrCreateUser
}
