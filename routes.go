package main

import (
	"linn221/shop/handlers"
	"linn221/shop/middlewares"
	"linn221/shop/models"
	"net/http"
)

func myRouter(c *Container) *http.ServeMux {

	mainMux := http.NewServeMux()
	authMux := http.NewServeMux()
	authMux.Handle("GET /me", handlers.Me(c.DB))
	authMux.HandleFunc("POST /change-password", handlers.ChangePassword(c.DB))
	authMux.HandleFunc("POST /update-profile", handlers.UpdateUserInfo(c.DB))

	//categories
	authMux.HandleFunc("POST /categories", handlers.HandleCategoryCreate(c.DB,
		c.Readers.CategoryListService.CleanCache,
	))
	authMux.HandleFunc("PUT /categories/{id}", handlers.HandleCategoryUpdate(c.DB,
		c.Readers.CategoryListService.CleanCache,
		c.Readers.CategoryGetService.CleanCache,
	))
	authMux.HandleFunc("DELETE /categories/{id}", handlers.HandleCategoryDelete(c.DB,
		c.Readers.CategoryListService.CleanCache,
		c.Readers.CategoryGetService.CleanCache,
	))
	authMux.HandleFunc("PATCH /categories/{id}/toggle", handlers.HandleToggleActive[models.Category](c.DB,
		c.Readers.CategoryGetService.CleanCache,
		c.Readers.CategoryListService.CleanCache,
	))
	authMux.HandleFunc("GET /categories/{id}", handlers.DefaultGetHandler(c.Readers.CategoryGetService))
	authMux.HandleFunc("GET /categories", handlers.DefaultListHandler(c.Readers.CategoryListService))

	mainMux.HandleFunc("POST /upload-single", handlers.HandleImageUploadSingle(c.DB, c.ImageDirectory))
	mainMux.HandleFunc("POST /register", handlers.Register(c.DB))
	mainMux.HandleFunc("POST /login", handlers.Login(c.DB, c.Cache))
	mainMux.Handle("/api/", http.StripPrefix("/api", middlewares.Auth(authMux)))
	return mainMux
}
