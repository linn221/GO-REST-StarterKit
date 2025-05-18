package handlers

import (
	"linn221/shop/myctx"
	"linn221/shop/services"
	"net/http"

	"gorm.io/gorm"
)

type GeneralHandlerFunc func(http.ResponseWriter, *http.Request, int, string, *gorm.DB, services.CacheService) error

func DefaultH(container *services.Container, h GeneralHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userId, shopId, err := myctx.GetIdsFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = h(w, r, userId, shopId, container.DB.WithContext(ctx), container.Cache)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}

type InputHandlerFunc[T any] func(http.ResponseWriter, *http.Request, *T, int, string, *gorm.DB, services.CacheService) error

func WithInput[T any](container *services.Container, h InputHandlerFunc[T]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userId, shopId, err := myctx.GetIdsFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		input, ok, err := parseJson[T](w, r, container.Validate)
		if !ok {
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		err = h(w, r, input, userId, shopId, container.DB.WithContext(ctx), container.Cache)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
