package middlewares

import (
	"fmt"
	"net/http"
)

func Recovery(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			r := recover()
			if r != nil {
				http.Error(w, fmt.Sprint(r), http.StatusInternalServerError)
			}
		}()
		h.ServeHTTP(w, r)
	})
}
