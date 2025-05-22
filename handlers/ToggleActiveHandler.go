package handlers

import (
	"linn221/shop/myctx"
	"linn221/shop/services"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

func HandleToggleActive[T services.HasIsActiveStatus](db *gorm.DB,
	cleanInstanceCache services.CleanInstanceCache,
	cleanListingCache services.CleanListingCache,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		_, shopId, err := myctx.GetIdsFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resIdStr := r.PathValue("id")
		resId, err := strconv.Atoi(resIdStr)
		if err != nil || resId <= 0 {
			http.Error(w, "invalid resource id", http.StatusBadRequest)
			return
		}

		var resource T
		if err := db.WithContext(ctx).Where("shop_id = ?", shopId).First(&resource, resId).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				finalErrHandle(w,
					respondNotFound(w, "record not found"),
				)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var isActive bool
		if r.URL.Query().Get("active") == "1" {
			isActive = true
		}
		if isActive != resource.GetIsActive() {

			tx := db.WithContext(ctx).Begin()
			defer tx.Rollback()

			if err := tx.Model(&resource).UpdateColumn("is_active", isActive).Error; err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if err := cleanInstanceCache(resId); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if err := cleanListingCache(shopId); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if err := tx.Commit().Error; err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		respondNoContent(w)
	}
}
