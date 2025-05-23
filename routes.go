package main

import (
	"linn221/shop/handlers"
	"linn221/shop/middlewares"
	"linn221/shop/models"
	"log"
	"net/http"
	"time"
)

func (app *App) Serve() {

	mainMux := http.NewServeMux()
	authMux := http.NewServeMux()
	authMux.Handle("GET /me", handlers.Me(app.DB))
	authMux.HandleFunc("POST /change-password", handlers.ChangePassword(app.DB))
	authMux.HandleFunc("POST /update-profile", handlers.UpdateUserInfo(app.DB))

	//categories
	authMux.HandleFunc("POST /categories", handlers.HandleCategoryCreate(app.DB,
		app.Readers.CategoryListService.CleanCache,
	))
	authMux.HandleFunc("PUT /categories/{id}", handlers.HandleCategoryUpdate(app.DB,
		app.Readers.CategoryListService.CleanCache,
		app.Readers.CategoryGetService.CleanCache,
	))
	authMux.HandleFunc("DELETE /categories/{id}", handlers.HandleCategoryDelete(app.DB,
		app.Readers.CategoryListService.CleanCache,
		app.Readers.CategoryGetService.CleanCache,
	))
	authMux.HandleFunc("PATCH /categories/{id}/toggle", handlers.HandleToggleActive[models.Category](app.DB,
		app.Readers.CategoryGetService.CleanCache,
		app.Readers.CategoryListService.CleanCache,
	))
	authMux.HandleFunc("GET /categories/{id}", handlers.DefaultGetHandler(app.Readers.CategoryGetService))
	authMux.HandleFunc("GET /categories", handlers.DefaultListHandler(app.Readers.CategoryListService))

	//units
	authMux.HandleFunc("POST /units", handlers.HandleUnitCreate(app.DB,
		app.Readers.UnitListService.CleanCache,
	))
	authMux.HandleFunc("PUT /units/{id}", handlers.HandleUnitUpdate(app.DB,
		app.Readers.UnitGetService.CleanCache,
		app.Readers.UnitListService.CleanCache,
	))
	authMux.HandleFunc("DELETE /units/{id}", handlers.HandleUnitDelete(app.DB,
		app.Readers.UnitGetService.CleanCache,
		app.Readers.UnitListService.CleanCache,
	))
	authMux.HandleFunc("PATCH /units/{id}/toggle", handlers.HandleToggleActive[models.Unit](app.DB,
		app.Readers.UnitGetService.CleanCache,
		app.Readers.UnitListService.CleanCache,
	))
	authMux.HandleFunc("GET /units/{id}", handlers.DefaultGetHandler(app.Readers.UnitGetService))
	authMux.HandleFunc("GET /units", handlers.DefaultListHandler(app.Readers.UnitListService))

	mainMux.HandleFunc("POST /upload-single", handlers.HandleImageUploadSingle(app.DB, app.ImageDirectory))
	mainMux.HandleFunc("POST /register", handlers.Register(app.DB))
	mainMux.HandleFunc("POST /login", handlers.Login(app.DB, app.Cache))

	mainMux.Handle("/api/", http.StripPrefix("/api", middlewares.Auth(app.ResourceRateLimit(authMux))))

	srv := http.Server{
		Addr:         ":" + app.Port,
		Handler:      middlewares.Default(app.GeneralRateLimit(mainMux), app.Cache),
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Fatal(srv.ListenAndServe())
}
