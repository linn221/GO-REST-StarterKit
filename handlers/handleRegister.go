package handlers

import (
	"encoding/json"
	"linn221/shop/models"
	"net/http"
)

type NewShop struct {
	Name     string `json:"name" validate:"required,min=4,max=255"`
	Email    string `json:"email" validate:"required,email,min=5,max=255"`
	Password string `json:"password" validate:"required,min=8,max=255"`
	PhoneNo  string `json:"phone_no" validate:"required,min=2"`
}

func Register(userService *models.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// var input NewShop
		// defer r.Body.Close()
		// err := json.NewDecoder(r.Body).Decode(&input)
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// }
		input, ok, finalErr := parseJson[NewShop](w, r)
		if !ok {
			if finalErr != nil {
				http.Error(w, finalErr.Error(), http.StatusInternalServerError)
			}
			return
		}

		ctx := r.Context()
		user, err := userService.Register(ctx, input.Name, input.Email, input.Password, input.PhoneNo)
		if err != nil {
			respondError(w, err)
			return
		}

		if err := json.NewEncoder(w).Encode(user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
