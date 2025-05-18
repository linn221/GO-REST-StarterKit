package handlers

import (
	"linn221/shop/services"
	"net/http"

	"gorm.io/gorm"
)

type MyInfo struct {
	Id       int    `json:"id"`
	ShopId   string `json:"shop_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	PhoneNo  string `json:"phone_no"`
}

func Me(w http.ResponseWriter, r *http.Request, userId int, shopId string, db *gorm.DB, cache services.CacheService) error {

	user, errs := services.UserService.GetUser(userId, db, cache)
	if errs != nil {
		return errs.Respond(w)
	}

	return respondOk(w, MyInfo{
		Id:       userId,
		ShopId:   shopId,
		Username: user.Username,
		Email:    user.Email,
		PhoneNo:  user.PhoneNo,
	})
}
