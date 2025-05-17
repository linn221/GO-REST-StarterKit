package main

import (
	"linn221/shop/handlers"
	"linn221/shop/services"
	"net/http"
)

func myRouter(container *services.Container) *http.ServeMux {

	mainMux := http.NewServeMux()
	authMux := http.NewServeMux()
	authMux.Handle("POST /users", handlers.DefaultH(container, nil))
	mainMux.Handle("/api/", http.StripPrefix("/api", authMux))
	return mainMux
}
