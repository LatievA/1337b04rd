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

func (h *UserHandler) UserRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/user", h.HandleSession)
}

func (h *UserHandler) HandleSession(w http.ResponseWriter, r *http.Request) {
	
}
