package main

import (
	"context"
	"linn221/shop/config"
	"linn221/shop/middlewares"
	"linn221/shop/services"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {

	godotenv.Load(".env")
	port := "8080"
	if p := os.Getenv("API_PORT"); p != "" {
		port = p
	}
	ctx := context.Background()
	var cacheService services.CacheService
	var db *gorm.DB
	cacheService = config.ConnectRedis(ctx)
	db = config.ConnectDB()
	container := &services.Container{
		DB:    db,
		Cache: cacheService,
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
