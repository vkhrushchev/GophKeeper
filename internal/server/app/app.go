package app

import (
	"github.com/gorilla/mux"
	"github.com/vkhrushchev/GophKeeper/internal/server/handlers"
	"github.com/vkhrushchev/GophKeeper/internal/server/usecase"
)

type App struct {
	r                   mux.Router
	registerUserUseCase *usecase.RegisterUserUseCase
	loginUseCase        *usecase.LoginUseCase
	logoutUseCase       *usecase.LogoutUseCase
}

func (a *App) RegisterRoutes() {
	a.r.Handle(
		"/api/v1/user/register",
		&handlers.RegisterUserHandler{},
	).Methods("POST")

	a.r.Handle(
		"/api/v1/auth/login",
		&handlers.LoginUserHandler{},
	).Methods("POST")

	a.r.Handle(
		"/api/v1/auth/logout",
		&handlers.LogoutUserHandler{},
	).Methods("POST")
}
