package handlers

import (
	"linn221/shop/models"
	"linn221/shop/services"
	"net/http"
)

type NewCategory struct {
	Name inputString `json:"name" validate:"required,min=2,max=100"`
}

func HandleCategoryCreate(categoryService services.CategoryCruder) http.HandlerFunc {
	return CreateHandler(func(w http.ResponseWriter, r *http.Request, session *CreateSession[NewCategory]) error {

		category := models.Category{Name: string(session.Input.Name)}
		cat, errs := categoryService.CreateCategory(r.Context(), &category)
		if errs != nil {
			return errs.Respond(w)
		}

		return respondOk(w, cat)
	})
}

func HandleCategoryUpdate(categoryService services.CategoryCruder) http.HandlerFunc {
	return UpdateHandler(func(w http.ResponseWriter, r *http.Request, session *UpdateSession[NewCategory]) error {
		// cat, errs := categoryService.UpdateCategory(r.Context(), session.ResId, session.Input)
		// if errs != nil {
		// 	return errs.Respond(w)
		// }

		// respondOk(w, cat)
		return nil

	})
}
