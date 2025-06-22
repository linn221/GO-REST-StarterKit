package handlers

import (
	"linn221/shop/models"
	"net/http"

	"gorm.io/gorm"
)

type NewCategory struct {
	Name        inputString     `json:"name" validate:"required,min=2,max=100"`
	Description *optionalString `json:"description" validate:"omitempty,max=1000"`
}

func HandleCategoryCreate(categoryService *models.CategoryService) http.HandlerFunc {
	return CreateHandler(func(w http.ResponseWriter, r *http.Request, session Session, input *NewCategory) error {
		ctx := r.Context()

		category := models.Category{
			Name:        input.Name.String(),
			Description: input.Description.StringPtr(),
		}
		category.ShopId = session.ShopId

		_, err := categoryService.Store(ctx, &category, session.ShopId)
		if err != nil {
			return err
		}

		w.WriteHeader(http.StatusCreated)
		return nil
	})
}

func HandleCategoryUpdate(db *gorm.DB,
	cleanCache func(db *gorm.DB, shopId string, id int) error,
) http.HandlerFunc {
	return UpdateHandler(func(w http.ResponseWriter, r *http.Request, session Session, input *NewCategory) error {
		ctx := r.Context()
		if errs := input.validate(db.WithContext(ctx), session.ShopId, session.ResId); errs != nil {
			return errs.Respond(w)
		}
		var category models.Category
		if err := db.WithContext(ctx).Where("shop_id = ?", session.ShopId).First(&category, session.ResId).Error; err != nil {
			return err
		}
		updates := map[string]any{
			"Name": input.Name.String(),
		}
		if input.Description.IsPresent() {
			updates["Description"] = input.Description
		}

		err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			if err := tx.Model(&category).Updates(updates).Error; err != nil {
				return err
			}
			if err := cleanCache(db.WithContext(ctx), session.ShopId, session.ResId); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}

		respondNoContent(w)
		return nil
	})
}

func HandleCategoryDelete(db *gorm.DB,
	cleanCache func(*gorm.DB, string, int) error,
) http.HandlerFunc {
	return DeleteHandler(func(w http.ResponseWriter, r *http.Request, session Session) error {
		ctx := r.Context()
		var category models.Category
		if err := db.WithContext(ctx).Where("shop_id = ?", session.ShopId).First(&category, session.ResId).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return respondNotFound(w, "category not found")
			}
			return err
		}

		if errs := Validate(db.WithContext(ctx),
			NewNoResultRule("items", "category has been used in items", NewFilter("category_id = ?", session.ResId)),
		); errs != nil {
			return errs.Respond(w)
		}
		err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			if err := tx.Delete(&category).Error; err != nil {
				return err
			}
			if err := cleanCache(db.WithContext(ctx), session.ShopId, session.ResId); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}

		respondNoContent(w)
		return nil
	})

}

// // func HandleCategoryGet(getService services.Getter[models.Category]) http.HandlerFunc {
// // 	return GetHandler(func(w http.ResponseWriter, r *http.Request, session *GetSession) error {
// // 		category, found, err := getService.Get(session.ShopId, session.ResId)
// // 		if err != nil {
// // 			return err
// // 		}
// // 		if !found {
// // 			return respondNotFound(w, "category not found")
// // 		}
// // 		return respondOk(w, category)
// // 	})
// // }

// func HandleCategoryList(listService services.Lister[models.Category]) http.HandlerFunc {
// 	return DefaultHandler(func(w http.ResponseWriter, r *http.Request, session *DefaultSession) error {
// 		categories, err := listService.List(session.ShopId)
// 		if err != nil {
// 			return err
// 		}
// 		return respondOk(w, categories)
// 	})
// }
