package main

import (
	"context"
	"linn221/shop/config"
	"linn221/shop/middlewares"
	"linn221/shop/models"
	"linn221/shop/services"
	"log"
	"net/http"
	"time"

	"gorm.io/gorm"
)

func main() {

	// godotenv.Load(".env")
	port := "8080"
	// if p := os.Getenv("API_PORT"); p != "" {
	// 	port = p
	// }
	ctx := context.Background()
	var db *gorm.DB
	db = config.ConnectDB()
	var cacheService services.CacheService
	cacheService = config.ConnectRedis(ctx)
	dir := config.GetImageDirectory()
	readServices := models.NewReadServices(db, cacheService)

	container := &Container{
		DB:             db,
		Cache:          cacheService,
		ImageDirectory: dir,
		ReadServices:   readServices,
	}

	mux := myRouter(container)
	srv := http.Server{
		Addr:         ":" + port,
		Handler:      middlewares.Default(mux, cacheService),
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Fatal(srv.ListenAndServe())
}
