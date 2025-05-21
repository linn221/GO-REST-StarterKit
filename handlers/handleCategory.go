package handlers

import (
	"linn221/shop/models"
	"linn221/shop/services"
	"net/http"

	"gorm.io/gorm"
)

type NewCategory struct {
	Name        inputString    `json:"name" validate:"required,min=2,max=100"`
	Description optionalString `json:"description" validate:"omitempty,max=1000"`
}

func (input *NewCategory) validate(db *gorm.DB, shopId string, id int) *ServiceError {
	return Validate(db, NewUniqueRule("categories", "name", input.Name.String(), id, "duplicate category name").
		Filter("shop_id = ?", shopId))
}

func HandleCategoryCreate(db *gorm.DB, cleanListingCache services.ListingCacheCleaner) http.HandlerFunc {
	return CreateHandler(func(w http.ResponseWriter, r *http.Request, session *CreateSession[NewCategory]) error {
		ctx := r.Context()
		if errs := session.Input.validate(db.WithContext(ctx), session.ShopId, 0); errs != nil {
			return errs.Respond(w)
		}

		category := models.Category{
			Name:        session.Input.Name.String(),
			Description: session.Input.Description.StringPtr(),
		}
		category.ShopId = session.ShopId

		if err := db.WithContext(ctx).Create(&category).Error; err != nil {
			return err
		}
		if err := cleanListingCache(session.ShopId); err != nil {
			return err
		}

		w.WriteHeader(http.StatusCreated)
		return nil
	})
}
func HandleCategoryUpdate(db *gorm.DB, cleanListingCache services.ListingCacheCleaner, cleanInstanceCache services.InstanceCacheCleaner) http.HandlerFunc {
	return UpdateHandler(func(w http.ResponseWriter, r *http.Request, session *UpdateSession[NewCategory]) error {
		ctx := r.Context()
		if errs := session.Input.validate(db.WithContext(ctx), session.ShopId, session.ResId); errs != nil {
			return errs.Respond(w)
		}
		var category models.Category
		if err := db.WithContext(ctx).Where("shop_id = ?", session.ShopId).First(&category, session.ResId).Error; err != nil {
			return err
		}
		updates := map[string]any{
			"Name": session.Input.Name.String(),
		}
		if session.Input.Description.IsPresent() {
			updates["Description"] = session.Input.Description
		}
		if err := db.WithContext(ctx).Model(&category).Updates(updates).Error; err != nil {
			return err
		}
		if err := cleanInstanceCache(session.ResId); err != nil {
			return err
		}
		if err := cleanListingCache(session.ShopId); err != nil {
			return err
		}

		respondNoContent(w)
		return nil
	})
}

func HandleCategoryDelete(db *gorm.DB) http.HandlerFunc {
	return DeleteHandler(func(w http.ResponseWriter, r *http.Request, session *DeleteSession) error {
		ctx := r.Context()
		var category models.Category
		if err := db.WithContext(ctx).Where("shop_id = ?", session.ShopId).First(&category, session.ResId).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return respondNotFound(w, "category not found")
			}
			return err
		}

		// if errs := Validate(db.WithContext(ctx),
		// 	NewNoResultRule("items", "category has been used in items", NewFilter("category_id = ?", session.ResId)),
		// ); errs != nil {
		// 	return errs.Respond(w)
		// }

		if err := db.WithContext(ctx).Delete(&category).Error; err != nil {
			return err
		}

		respondNoContent(w)
		return nil
	})

}

func HandleCategoryGet(getService services.Getter[models.Category]) http.HandlerFunc {
	return GetHandler(func(w http.ResponseWriter, r *http.Request, session *GetSession) error {
		category, found, err := getService.Get(session.ShopId, session.ResId)
		if err != nil {
			return err
		}
		if !found {
			return respondNotFound(w, "category not found")
		}
		return respondOk(w, category)
	})
}

func HandleCategoryList(db *gorm.DB) http.HandlerFunc {
	return DefaultHandler(func(w http.ResponseWriter, r *http.Request, session *DefaultSession) error {
		var categories []models.Category
		if err := db.WithContext(r.Context()).Where("shop_id = ?", session.ShopId).Find(&categories).Error; err != nil {
			return err
		}
		return respondOk(w, categories)
	})
}
