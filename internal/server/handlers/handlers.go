package handlers

import (
	"github.com/vkhrushchev/GophKeeper/internal/server/usecase"
	"net/http"
)

type RegisterUserHandler struct {
	registerUserUseCase *usecase.RegisterUserUseCase
}

func (h *RegisterUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

type LoginUserHandler struct {
	loginUseCase *usecase.LoginUseCase
}

func (h *LoginUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

type LogoutUserHandler struct {
	logoutUseCase *usecase.LogoutUseCase
}

func (h *LogoutUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
