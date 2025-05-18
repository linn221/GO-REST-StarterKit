package handlers

import (
	"linn221/shop/models"
	"linn221/shop/services"
	"linn221/shop/utils"
	"net/http"

	"gorm.io/gorm"
)

type LoginInfo struct {
	Username string `json:"username" validate:"required,max=255"`
	Password string `json:"password" validate:"required,max=255"`
}

func Login(container *services.Container) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		login, ok, err := parseJson[LoginInfo](w, r, container.Validate)
		if !ok {
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		var user models.User
		if err := container.DB.Where("username = ?", login.Username).First(&user).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				http.Error(w, "invalid username/password", http.StatusBadRequest)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := utils.ComparePassword(user.Password, login.Password); err != nil {
			http.Error(w, "invalid username/password", http.StatusBadRequest)
			return
		}
		token, err := services.NewSession(user.Id, user.ShopId, container.Cache)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		respondOk(w, map[string]string{
			"token": token,
		})
	}
}
