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

type CreateFunc[T any] func(w http.ResponseWriter, r *http.Request, session *CreateSession[T]) error

func CreateHandler[T any](handle CreateFunc[T]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		userId, shopId, err := myctx.GetIdsFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		input, ok, err := parseJson[T](w, r)
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
