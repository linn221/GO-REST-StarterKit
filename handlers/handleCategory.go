package handlers

import (
	"linn221/shop/models"
	"net/http"
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

func HandleCategoryUpdate(categoryService *models.CategoryService) http.HandlerFunc {
	return UpdateHandler(func(w http.ResponseWriter, r *http.Request, session Session, input *NewCategory) error {
		ctx := r.Context()

		_, err := categoryService.Update(ctx, &models.Category{
			Name:        input.Name.String(),
			Description: input.Description.StringPtr(),
		}, session.ResId, session.ShopId)
		if err != nil {
			return err
		}
		respondNoContent(w)
		return nil
	})
}

func HandleCategoryDelete(categoryService *models.CategoryService) http.HandlerFunc {
	return DeleteHandler(func(w http.ResponseWriter, r *http.Request, session Session) error {
		ctx := r.Context()
		_, err := categoryService.Delete(ctx, session.ResId, session.ShopId)
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
