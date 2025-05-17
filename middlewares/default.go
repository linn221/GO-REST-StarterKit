package middlewares

import (
	"linn221/shop/services"
	"net/http"
)

func Default(h http.Handler, cache services.CacheService) http.Handler {
	sessionMd := SessionMiddleware{Cache: cache}

	return Recovery(sessionMd.Middleware(Logging(h)))
}
