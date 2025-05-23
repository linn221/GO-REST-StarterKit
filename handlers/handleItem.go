package handlers

import (
	"linn221/shop/models"
	"linn221/shop/services"
	"net/http"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type NewItem struct {
	Name          inputString     `json:"name" validate:"required,min=2,max=100"`
	Description   *optionalString `json:"description" validate:"omitempty,max=500"`
	CategoryId    int             `json:"category_id" validate:"required,number,gte=1"`
	UnitId        int             `json:"unit_id" validate:"required,number,gte=1"`
	SalesPrice    decimal.Decimal `json:"sales_price" validate:"required,number"`
	PurchasePrice decimal.Decimal `json:"purchase_price" validate:"requried,number"`
}

func (input *NewItem) validate(db *gorm.DB, shopId string, id int) *ServiceError {
	panic("not implemented") //2d
}

func HandleItemCreate(db *gorm.DB,
	cleanListingCache services.CleanListingCache,
) http.HandlerFunc {
	return CreateHandler(func(w http.ResponseWriter, r *http.Request, session Session, input *NewItem) error {

		ctx := r.Context()
		if errs := input.validate(db.WithContext(ctx), session.ShopId, 0); errs != nil {
			return errs.Respond(w)
		}
		item := models.Item{
			Name:          input.Name.String(),
			Description:   input.Description.StringPtr(),
			CategoryId:    input.CategoryId,
			UnitId:        input.UnitId,
			SalesPrice:    input.SalesPrice,
			PurchasePrice: input.PurchasePrice,
		}
		item.ShopId = session.ShopId

		err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(&item).Error; err != nil {
				return err
			}

			if err := cleanListingCache(session.ShopId); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}

		w.WriteHeader(http.StatusCreated)
		return nil
	})
}

func HandleItemUpdate(db *gorm.DB,
	cleanInstanceCache services.CleanInstanceCache,
	cleanListingCache services.CleanListingCache,
) http.HandlerFunc {
	return UpdateHandler(func(w http.ResponseWriter, r *http.Request, session Session, input *NewItem) error {

		ctx := r.Context()
		if errs := input.validate(db.WithContext(ctx), session.ShopId, session.ResId); errs != nil {
			return errs.Respond(w)
		}
		item, errs := first[models.Item](db.WithContext(ctx), session.ShopId, session.ResId)
		if errs != nil {
			return errs.Respond(w)
		}

		updates := map[string]any{
			"Name":          input.Name,
			"CategoryId":    input.CategoryId,
			"UnitId":        input.UnitId,
			"PurchasePrice": input.PurchasePrice,
			"SalesPrice":    input.SalesPrice,
		}
		if input.Description.IsPresent() {
			updates["Description"] = input.Description.String()
		}

		err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			if err := tx.Model(&item).Updates(updates).Error; err != nil {
				return err
			}
			if err := cleanInstanceCache(session.ResId); err != nil {
				return err
			}
			if err := cleanListingCache(session.ShopId); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}

		respondNoContent(w)
		return nil
	})
}

func HandleItemDelete(db *gorm.DB,
	cleanInstanceCache services.CleanInstanceCache,
	cleanListingCache services.CleanListingCache,
) http.HandlerFunc {
	return DeleteHandler(func(w http.ResponseWriter, r *http.Request, session Session) error {

		ctx := r.Context()
		item, errs := first[models.Item](db.WithContext(ctx), session.ShopId, session.ResId)
		if errs != nil {
			return errs.Respond(w)
		}

		// valiate
		err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			if err := tx.Delete(&item).Error; err != nil {
				return err
			}
			if err := cleanInstanceCache(session.ResId); err != nil {
				return err
			}
			if err := cleanListingCache(session.ShopId); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}

		respondNoContent(w)
		return nil
	})
}
