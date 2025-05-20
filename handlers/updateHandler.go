package handlers

import (
	"linn221/shop/myctx"
	"net/http"
	"strconv"
)

type UpdateSession[T any] struct {
	UserId int
	ShopId string
	ResId  int
	Input  *T
}

type UpdateFunc[T any] func(http.ResponseWriter, *http.Request, *UpdateSession[T]) error

func UpdateHandler[T any](handle UpdateFunc[T]) http.HandlerFunc {
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

		resIdStr := r.PathValue("id")
		resId, err := strconv.Atoi(resIdStr)
		if err != nil || resId <= 0 {
			http.Error(w, "invalid resource id", http.StatusBadRequest)
			return
		}

		UpdateSession := UpdateSession[T]{
			UserId: userId,
			ShopId: shopId,
			Input:  input,
			ResId:  resId,
		}
		err = handle(w, r, &UpdateSession)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
