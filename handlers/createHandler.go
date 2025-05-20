package handlers

import (
	"linn221/shop/myctx"
	"net/http"
)

type CreateSession[T any] struct {
	UserId int
	ShopId string
	Input  *T
}

func CreateHandler[T any](handle func(http.ResponseWriter, *http.Request, *CreateSession[T]) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		userId, shopId, err := myctx.GetIdsFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		input, ok, err := parseJson2[T](w, r)
		if !ok {
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		CreateSession := CreateSession[T]{
			UserId: userId,
			ShopId: shopId,
			Input:  input,
		}
		err = handle(w, r, &CreateSession)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
