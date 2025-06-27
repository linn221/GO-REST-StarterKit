package handlers

import (
	"linn221/shop/models"
	"net/http"
	"strconv"

	"github.com/shopspring/decimal"
)

type NewItem struct {
	Name          inputString      `json:"name" validate:"required,min=2,max=100"`
	Description   *optionalString  `json:"description" validate:"omitempty,max=500"`
	CategoryId    int              `json:"category_id" validate:"required,number,gte=1"`
	UnitId        int              `json:"unit_id" validate:"required,number,gte=1"`
	SalesPrice    *decimal.Decimal `json:"sales_price" validate:"required,number,gte=1"`
	PurchasePrice *decimal.Decimal `json:"purchase_price" validate:"required,number,gte=1"`
}

func HandleItemCreate(itemService *models.ItemService) http.HandlerFunc {
	return CreateHandler(func(w http.ResponseWriter, r *http.Request, session Session, input *NewItem) error {

		ctx := r.Context()
		item := models.Item{
			Name:          input.Name.String(),
			Description:   input.Description.StringPtr(),
			CategoryId:    input.CategoryId,
			UnitId:        input.UnitId,
			SalesPrice:    *input.SalesPrice,
			PurchasePrice: *input.PurchasePrice,
		}
		item.ShopId = session.ShopId

		_, err := itemService.Store(ctx, &item, session.ShopId)
		if err != nil {
			return err
		}

		w.WriteHeader(http.StatusCreated)
		return nil
	})
}

func HandleItemUpdate(itemService *models.ItemService) http.HandlerFunc {
	return UpdateHandler(func(w http.ResponseWriter, r *http.Request, session Session, input *NewItem) error {

		ctx := r.Context()

		updates := models.Item{
			Name:          input.Name.String(),
			CategoryId:    input.CategoryId,
			UnitId:        input.UnitId,
			PurchasePrice: *input.PurchasePrice,
			SalesPrice:    *input.SalesPrice,
			Description:   input.Description.StringPtr(),
		}

		_, err := itemService.Update(ctx, &updates, session.ResId, session.ShopId)
		if err != nil {
			return err
		}

		respondNoContent(w)
		return nil
	})
}

func HandleItemDelete(itemService *models.ItemService) http.HandlerFunc {
	return DeleteHandler(func(w http.ResponseWriter, r *http.Request, session Session) error {

		ctx := r.Context()
		_, err := itemService.Delete(ctx, session.ResId, session.ShopId)
		if err != nil {
			return err
		}

		respondNoContent(w)
		return nil
	})
}

func parseItemSearch(r *http.Request) (*models.ItemSearch, error) {
	var search models.ItemSearch
	var err error
	if s := r.URL.Query().Get("search"); s != "" {
		search.Search = s
	}
	if s := r.URL.Query().Get("category_id"); s != "" {
		search.CategoryId, err = strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
	}
	if s := r.URL.Query().Get("unit_id"); s != "" {
		search.UnitId, err = strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
	}

	if min := r.URL.Query().Get("sales_price_min"); min != "" {
		salesPriceMin, err := decimal.NewFromString(min)
		if err != nil {
			return nil, err
		}
		search.SalesPriceMin = &salesPriceMin
		max := r.URL.Query().Get("sales_price_max")
		if max != "" {
			salesPriceMax, err := decimal.NewFromString(max)
			if err != nil {
				return nil, err
			}
			search.SalesPriceMax = &salesPriceMax
		}
	}

	if min := r.URL.Query().Get("purchase_price_min"); min != "" {
		purchasePriceMin, err := decimal.NewFromString(min)
		if err != nil {
			return nil, err
		}
		search.PurchasePriceMin = &purchasePriceMin
		max := r.URL.Query().Get("purchase_price_max")
		if max != "" {
			purchasePriceMax, err := decimal.NewFromString(max)
			if err != nil {
				return nil, err
			}
			search.PurchasePriceMax = &purchasePriceMax
		}
	}
	return &search, nil
}

func HandleItemSearch(itemService *models.ItemService) http.HandlerFunc {
	return DefaultHandler(func(w http.ResponseWriter, r *http.Request, session *DefaultSession) error {
		search, err := parseItemSearch(r)
		if err != nil {
			return respondClientError(w, err.Error())
		}

		//2d cache results
		results, err := itemService.SearchItems(r.Context(), session.ShopId, search)
		if err != nil {
			return err
		}
		respondOk(w, results)
		return nil
	})
}
