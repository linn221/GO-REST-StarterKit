package handlers

import (
	"linn221/shop/myctx"
	"linn221/shop/services"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type ResourceHandlerFunc func(http.ResponseWriter, *http.Request, int, int, int, *gorm.DB, services.CacheService) error

func ResourceH(container *services.Container, h ResourceHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		userId, shopId, err := myctx.GetIdsFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		resIdStr := r.PathValue("id")
		if resIdStr == "" {
			http.Error(w, "id param is required", http.StatusBadRequest)
			return
		}
		resId, err := strconv.Atoi(resIdStr)
		if err != nil {
			http.Error(w, "id param is required", http.StatusBadRequest)
			return
		}

		err = h(w, r, userId, shopId, resId, container.DB.WithContext(ctx), container.Cache)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}
