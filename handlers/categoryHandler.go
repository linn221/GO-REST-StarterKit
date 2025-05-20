package handlers

import (
	"linn221/shop/models"
	"linn221/shop/services"
	"net/http"

	"gorm.io/gorm"
)

func HandleCategoryCreate(db *gorm.DB, cache services.CacheService, categoryService services.CategoryCruder) http.HandlerFunc {
	return CreateHandler(func(w http.ResponseWriter, r *http.Request, cs *CreateSession[models.Category]) error {
		cat, errs := categoryService.CreateCategory(cs.Input, db, cache)
		if errs != nil {
			return errs.Respond(w)
		}

		respondOk(w, cat)
		return nil
	})
}
