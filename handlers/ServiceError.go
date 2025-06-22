package handlers

import (
	"linn221/shop/services"
	"net/http"
)

func ServiceErrorHandler(h func(http.ResponseWriter, *http.Request) *services.ServiceError) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serr := h(w, r)
		if serr != nil {
			if err := serr.Respond(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}
