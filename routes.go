package main

import (
	"fmt"
	"linn221/shop/handlers"
	"linn221/shop/middlewares"
	"linn221/shop/models"
	"linn221/shop/myctx"
	"linn221/shop/utils"
	"net/http"
	"time"
)

func myRouter(c *Container) http.Handler {

	mainMux := http.NewServeMux()
	authMux := http.NewServeMux()
	authMux.Handle("GET /me", handlers.Me(c.DB))
	authMux.HandleFunc("POST /change-password", handlers.ChangePassword(c.DB))
	authMux.HandleFunc("POST /update-profile", handlers.UpdateUserInfo(c.DB))

	// rate Limit
	// rate limiting crud endpoints by userId
	resourceRateLimit := middlewares.NewRateLimiter(c.Cache.GetClient(), time.Minute*5, 2000, "r", func(r *http.Request) (string, error) {
		ctx := r.Context()
		userId, _, err := myctx.GetIdsFromContext(ctx)
		if err != nil {
			return "", err
		}
		return fmt.Sprint(userId), nil
	})
	// rate limit by IP address for all routes
	generalRateLimit := middlewares.NewRateLimiter(c.Cache.GetClient(), time.Minute*2, 300, "g", func(r *http.Request) (string, error) {
		ip := utils.GetClientIP(r)
		return ip, nil
	})

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

	//units
	authMux.HandleFunc("POST /units", handlers.HandleUnitCreate(c.DB,
		c.Readers.UnitListService.CleanCache,
	))
	authMux.HandleFunc("PUT /units/{id}", handlers.HandleUnitUpdate(c.DB,
		c.Readers.UnitGetService.CleanCache,
		c.Readers.UnitListService.CleanCache,
	))
	authMux.HandleFunc("DELETE /units/{id}", handlers.HandleUnitDelete(c.DB,
		c.Readers.UnitGetService.CleanCache,
		c.Readers.UnitListService.CleanCache,
	))
	authMux.HandleFunc("PATCH /units/{id}/toggle", handlers.HandleToggleActive[models.Unit](c.DB,
		c.Readers.UnitGetService.CleanCache,
		c.Readers.UnitListService.CleanCache,
	))
	authMux.HandleFunc("GET /units/{id}", handlers.DefaultGetHandler(c.Readers.UnitGetService))
	authMux.HandleFunc("GET /units", handlers.DefaultListHandler(c.Readers.UnitListService))

	mainMux.HandleFunc("POST /upload-single", handlers.HandleImageUploadSingle(c.DB, c.ImageDirectory))
	mainMux.HandleFunc("POST /register", handlers.Register(c.DB))
	mainMux.HandleFunc("POST /login", handlers.Login(c.DB, c.Cache))
	mainMux.Handle("/api/", http.StripPrefix("/api", middlewares.Auth(resourceRateLimit(authMux))))

	return generalRateLimit(mainMux)
}
