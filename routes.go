package main

import (
	"linn221/shop/handlers"
	"linn221/shop/middlewares"
	"linn221/shop/services"
	"net/http"
)

func myRouter(container *services.Container) *http.ServeMux {

	mainMux := http.NewServeMux()
	authMux := http.NewServeMux()
	authMux.Handle("GET /me", handlers.Me(container.MyServices.UserService))
	authMux.HandleFunc("POST /change-password", handlers.ChangePassword(container.MyServices.UserService))
	authMux.HandleFunc("POST /update-profile", handlers.UpdateUserInfo(container.MyServices.UserService))
	mainMux.HandleFunc("POST /upload-single", handlers.HandleImageUploadSingle(container.MyServices.ImageService))
	mainMux.HandleFunc("POST /register", handlers.Register(container))
	mainMux.HandleFunc("POST /login", handlers.Login(container))
	mainMux.Handle("/api/", http.StripPrefix("/api", middlewares.Auth(authMux)))
	return mainMux
}
