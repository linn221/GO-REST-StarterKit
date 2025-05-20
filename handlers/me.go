package handlers

import (
	"linn221/shop/services"
	"net/http"
)

type MyInfo struct {
	Id       int    `json:"id"`
	ShopId   string `json:"shop_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	PhoneNo  string `json:"phone_no"`
}

func Me(userService services.UserCruder) http.HandlerFunc {

	return DefaultHandler(func(ds *DefaultSession, w http.ResponseWriter, r *http.Request) error {
		user, errs := userService.GetUser(r.Context(), ds.UserId)
		if errs != nil {
			return errs.Respond(w)
		}
		return respondOk(w, MyInfo{
			Id:       ds.UserId,
			ShopId:   ds.ShopId,
			Username: user.Username,
			Email:    user.Email,
			PhoneNo:  user.PhoneNo,
		})
	})
}
