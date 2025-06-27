package handlers

import (
	"linn221/shop/models"
	"linn221/shop/services"
	"net/http"
)

type LoginInfo struct {
	Username string `json:"username" validate:"required,max=255"`
	Password string `json:"password" validate:"required,max=255"`
}

func Login(userService *models.UserService, cache services.CacheService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		login, ok, err := parseJson[LoginInfo](w, r)
		if !ok {
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		user, err := userService.Login(ctx, login.Username, login.Password)
		if err != nil {
			respondError(w, err)
			return
		}

		token, err := services.NewSession(user.Id, user.ShopId, cache)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		respondOk(w, map[string]string{
			"token": token,
		})
	}
}
func Logout(cache services.CacheService) http.HandlerFunc {
	return DefaultHandler(func(w http.ResponseWriter, r *http.Request, session *DefaultSession) error {
		token := r.Header.Get("Token")
		if err := services.RemoveSession(token, cache); err != nil {
			return err
		}

		respondNoContent(w)
		return nil
	})
}
