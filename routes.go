package main

import (
	"linn221/shop/handlers"
	"linn221/shop/middlewares"
	"log"
	"net/http"
	"time"
)

func (app *App) Serve() {

	authMux := http.NewServeMux()

	authMux.Handle("GET /me", handlers.Me(app.DB))
	authMux.HandleFunc("POST /change-password", handlers.ChangePassword(app.Services.UserService))
	authMux.HandleFunc("POST /logout", handlers.Logout(app.Cache))
	authMux.HandleFunc("POST /update-profile", handlers.UpdateUserInfo(app.Services.UserService))

	//categories
	authMux.HandleFunc("POST /categories", handlers.HandleCategoryCreate(app.Services.CategoryService))
	authMux.HandleFunc("PUT /categories/{id}", handlers.HandleCategoryUpdate(app.Services.CategoryService))
	authMux.HandleFunc("DELETE /categories/{id}", handlers.HandleCategoryDelete(app.Services.CategoryService))
	// authMux.HandleFunc("PATCH /categories/{id}/toggle",
	// 	handlers.HandleToggleActive[models.Category](app.DB,
	// 		app.Readers.AfterCategoryUpdate,
	// 	))
	authMux.HandleFunc("GET /categories/{id}", handlers.DefaultGetHandler(app.Services.CategoryService.Get))
	authMux.HandleFunc("GET /categories", handlers.DefaultListHandler(app.Services.CategoryService.List))
	// authMux.HandleFunc("GET /categories/inactive",
	// 	handlers.ListInactiveHandler[models.Category, models.CategoryResource](app.DB),
	// )

	//units
	authMux.HandleFunc("POST /units", handlers.HandleUnitCreate(app.Services.UnitService))
	authMux.HandleFunc("PUT /units/{id}", handlers.HandleUnitUpdate(app.Services.UnitService))
	authMux.HandleFunc("DELETE /units/{id}", handlers.HandleUnitDelete(app.Services.UnitService))
	// authMux.HandleFunc("PATCH /units/{id}/toggle", handlers.HandleToggleActive[models.Unit](app.DB,
	// 	app.Readers.AfterUnitUpdate,
	// ))
	authMux.HandleFunc("GET /units/{id}", handlers.DefaultGetHandler(app.Services.UnitService.Get))
	authMux.HandleFunc("GET /units", handlers.DefaultListHandler(app.Services.UnitService.List))
	// authMux.HandleFunc("GET /units/inactive", handlers.ListInactiveHandler[models.Unit, models.UnitResource](app.DB))

	// items
	authMux.HandleFunc("POST /items", handlers.HandleItemCreate(app.Services.ItemService))
	authMux.HandleFunc("PUT /items/{id}", handlers.HandleItemUpdate(app.Services.ItemService))
	authMux.HandleFunc("DELETE /items/{id}", handlers.HandleItemDelete(app.Services.ItemService))
	// authMux.HandleFunc("PATCH /items/{id}/toggle",
	// 	handlers.HandleToggleActive[models.Item](app.DB, app.Readers.AfterUpdateItem),
	// )
	authMux.HandleFunc("GET /items/{id}", handlers.DefaultGetHandler(app.Services.ItemService.Get))
	authMux.HandleFunc("GET /items", handlers.DefaultListHandler(app.Services.ItemService.List))
	authMux.HandleFunc("GET /items/search", handlers.HandleItemSearch(app.Services.ItemService))
	// authMux.HandleFunc("GET /items/inactive",
	// 	handlers.ListCustomInactiveHandler(app.DB, models.FetchInactiveItemResources),
	// )

	mainMux := http.NewServeMux()
	// public routes
	mainMux.HandleFunc("POST /upload-single", handlers.HandleImageUploadSingle(app.DB, app.ImageDirectory))
	mainMux.HandleFunc("POST /register", handlers.Register(app.Services.UserService))
	mainMux.HandleFunc("POST /login", handlers.Login(app.Services.UserService, app.Cache))

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
