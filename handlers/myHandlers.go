package handlers

import (
	"linn221/shop/myctx"
	"linn221/shop/services"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type CreateFunc[T any] func(w http.ResponseWriter, r *http.Request, session Session, input *T) error

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

		CreateSession := Session{
			UserId: userId,
			ShopId: shopId,
		}
		err = handle(w, r, CreateSession, input)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

type DefaultSession struct {
	UserId int
	ShopId string
}

func DefaultHandler(handle func(http.ResponseWriter, *http.Request, *DefaultSession) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userId, shopId, err := myctx.GetIdsFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		session := DefaultSession{
			UserId: userId,
			ShopId: shopId,
		}

		err = handle(w, r, &session)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func DefaultListHandler[T any](listService services.Lister[T]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		_, shopId, err := myctx.GetIdsFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resourceSlice, err := listService.List(shopId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		finalErrHandle(w,
			respondOk(w, resourceSlice),
		)
	}
}

type Session struct {
	UserId int
	ShopId string
	ResId  int
}

type DeleteFunc func(w http.ResponseWriter, r *http.Request, session Session) error

func DeleteHandler(handle DeleteFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		userId, shopId, err := myctx.GetIdsFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resIdStr := r.PathValue("id")
		resId, err := strconv.Atoi(resIdStr)
		if err != nil || resId <= 0 {
			http.Error(w, "invalid resource id", http.StatusBadRequest)
			return
		}

		DeleteSession := Session{
			UserId: userId,
			ShopId: shopId,
			ResId:  resId,
		}
		err = handle(w, r, DeleteSession)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

type GetSession struct {
	UserId int
	ShopId string
	ResId  int
}

type GetFunc func(w http.ResponseWriter, r *http.Request, session *GetSession) error

func DefaultGetHandler[T any](getService services.Getter[T]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		_, shopId, err := myctx.GetIdsFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resIdStr := r.PathValue("id")
		resId, err := strconv.Atoi(resIdStr)
		if err != nil || resId <= 0 {
			http.Error(w, "invalid resource id", http.StatusBadRequest)
			return
		}

		resource, found, err := getService.Get(shopId, resId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !found {
			finalErrHandle(w,
				respondNotFound(w, "record not found"),
			)
			return
		}

		finalErrHandle(w,
			respondOk(w, resource),
		)
	}
}

func finalErrHandle(w http.ResponseWriter, err error) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func InputHandler[T any](handle func(http.ResponseWriter, *http.Request, *DefaultSession, *T) error) http.HandlerFunc {
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

		session := DefaultSession{
			UserId: userId,
			ShopId: shopId,
		}
		err = handle(w, r, &session, input)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// using gorm's smart scan
func ListInactiveHandler[Model any, Resource any](db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, shopId, err := myctx.GetIdsFromContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var results []Resource
		var m Model
		err = db.WithContext(r.Context()).Model(&m).Where("is_active = 0 AND shop_id = ?", shopId).Find(&results).Error
		finalErrHandle(w, err)

		finalErrHandle(w,
			respondOk(w, results),
		)
	}
}

func ListCustomInactiveHandler[Resource any](db *gorm.DB, fetch func(*gorm.DB, string) ([]Resource, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, shopId, err := myctx.GetIdsFromContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		results, err := fetch(db, shopId)
		finalErrHandle(w, err)

		finalErrHandle(w,
			respondOk(w, results),
		)
	}
}

func HandleToggleActive[T services.HasIsActiveStatus](db *gorm.DB,
	cleanCache func(*gorm.DB, string, int) error,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		_, shopId, err := myctx.GetIdsFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resIdStr := r.PathValue("id")
		resId, err := strconv.Atoi(resIdStr)
		if err != nil || resId <= 0 {
			http.Error(w, "invalid resource id", http.StatusBadRequest)
			return
		}

		var resource T
		if err := db.WithContext(ctx).Where("shop_id = ?", shopId).First(&resource, resId).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				finalErrHandle(w,
					respondNotFound(w, "record not found"),
				)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var isActive bool
		if r.URL.Query().Get("active") == "1" {
			isActive = true
		}
		if isActive != resource.GetIsActive() {

			tx := db.WithContext(ctx).Begin()
			defer tx.Rollback()

			if err := tx.Model(&resource).UpdateColumn("is_active", isActive).Error; err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if err := cleanCache(db.WithContext(ctx), shopId, resId); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if err := tx.Commit().Error; err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		respondNoContent(w)
	}
}

type UpdateFunc[T any] func(w http.ResponseWriter, r *http.Request, session Session, input *T) error

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

		UpdateSession := Session{
			UserId: userId,
			ShopId: shopId,
			ResId:  resId,
		}
		err = handle(w, r, UpdateSession, input)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
