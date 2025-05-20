package handlers

import (
	"linn221/shop/services"
	"net/http"
)

type NewPassword struct {
	OldPassword string `json:"old_password" validate:"required,max=255"`
	NewPassword string `json:"new_password" validate:"required,max=255,min=8"`
}

func ChangePassword(userSerivce services.UserCruder) http.HandlerFunc {
	return InputHandler[NewPassword](func(ds *DefaultSession, np *NewPassword, w http.ResponseWriter, r *http.Request) error {
		errs := userSerivce.ChangePassword(r.Context(), ds.UserId, np.OldPassword, np.NewPassword)
		if errs != nil {
			errs.Respond(w)
		}
		respondNoContent(w)
		return nil
	})
}

type NewUserEdit struct {
	Username inputString    `json:"username" validate:"required,min=3,max=100"`
	Email    inputString    `json:"email" validate:"required,email,min=4,max=100"`
	PhoneNo  optionalString `json:"phone_no" validate:"omitempty,min=5,max=20"`
}

func UpdateUserInfo(userSerivce services.UserCruder) http.HandlerFunc {
	return InputHandler[NewUserEdit](func(ds *DefaultSession, t *NewUserEdit, w http.ResponseWriter, r *http.Request) error {
		errs := userSerivce.UpdateInfo(r.Context(), ds.UserId, string(t.Username), string(t.Email), t.PhoneNo.StringPtr())
		if errs != nil {
			return errs.Respond(w)
		}

		respondNoContent(w)
		return nil
	})
}

// func CreateItemCategory(db *gorm.DB, cache services.CacheService, categoryService services.CategoryCruder, validate *validator.Validate) http.HandlerFunc {
// 	return inputHandler(func(w http.ResponseWriter, r *http.Request, input *models.Category, userId int, shopId string) error {
// 		cat, err := categoryService.CreateCategory(input, db, cache)
// 		if err != nil {
// 			return err
// 		}

// 		return respondOk(w, cat)
// 	}, validate)
// }

// func inputHandler[T any](f func(w http.ResponseWriter, r *http.Request, input *T, userId int, shopId string) error, validate *validator.Validate) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		ctx := r.Context()
// 		userId, shopId, err := myctx.GetIdsFromContext(ctx)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}

// 		input, ok, err := parseJson[T](w, r, validate)
// 		if !ok {
// 			if err != nil {
// 				http.Error(w, err.Error(), http.StatusInternalServerError)
// 			}
// 			return
// 		}

// 		err = f(w, r, input, userId, shopId)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 		}
// 	}
// }
