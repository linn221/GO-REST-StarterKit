package handlers

import (
	"linn221/shop/services"
	"net/http"

	"gorm.io/gorm"
)

type NewPassword struct {
	OldPassword string `json:"old_password" validate:"required,max=255"`
	NewPassword string `json:"new_password" validate:"required,max=255,min=8"`
}

func ChangePassword(w http.ResponseWriter, r *http.Request, input *NewPassword, userId int, shopId string, db *gorm.DB, cache services.CacheService) error {

	errs := services.UserService.ChangePassword(userId, input.OldPassword, input.NewPassword, db, cache)
	if errs != nil {
		return errs.Respond(w)
	}
	respondNoContent(w)
	return nil
}

type NewUserEdit struct {
	Name     inputString    `json:"name" validate:"required,min=2,max=100"`
	Username inputString    `json:"username" validate:"required,min=3,max=100"`
	Email    inputString    `json:"email" validate:"required,email,min=4,max=100"`
	PhoneNo  optionalString `json:"phone_no" validate:"omitempty,min=5,max=20"`
}

func UpdateUserInfo(w http.ResponseWriter, r *http.Request, input *NewUserEdit, userId int, shopId string, db *gorm.DB, cache services.CacheService) error {

	errs := services.UserService.UpdateInfo(userId, string(input.Name), string(input.Username), string(input.Email), input.PhoneNo.StringPtr(), db, cache)
	if errs != nil {
		return errs.Respond(w)
	}
	respondNoContent(w)
	return nil
}
