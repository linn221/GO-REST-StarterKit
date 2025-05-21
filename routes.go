package main

import (
	"linn221/shop/handlers"
	"linn221/shop/middlewares"
	"net/http"
)

func myRouter(container *Container) *http.ServeMux {

	mainMux := http.NewServeMux()
	authMux := http.NewServeMux()
	authMux.Handle("GET /me", handlers.Me(container.DB))
	authMux.HandleFunc("POST /change-password", handlers.ChangePassword(container.DB))
	authMux.HandleFunc("POST /update-profile", handlers.UpdateUserInfo(container.DB))

	//categories
	authMux.HandleFunc("POST /categories", handlers.HandleCategoryCreate(container.DB))
	authMux.HandleFunc("PUT /categories/{id}", handlers.HandleCategoryUpdate(container.DB, container.ReadServices.CategoryGetService.Clean))
	authMux.HandleFunc("DELETE /categories/{id}", handlers.HandleCategoryDelete(container.DB))
	authMux.HandleFunc("GET /categories/{id}", handlers.HandleCategoryGet(container.ReadServices.CategoryGetService))
	authMux.HandleFunc("GET /categories", handlers.HandleCategoryList(container.DB))

	mainMux.HandleFunc("POST /upload-single", handlers.HandleImageUploadSingle(container.DB, container.ImageDirectory))
	mainMux.HandleFunc("POST /register", handlers.Register(container.DB))
	mainMux.HandleFunc("POST /login", handlers.Login(container.DB, container.Cache))
	mainMux.Handle("/api/", http.StripPrefix("/api", middlewares.Auth(authMux)))
	return mainMux
}
