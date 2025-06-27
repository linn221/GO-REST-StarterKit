package models

import (
	"context"
	"linn221/shop/services"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Item struct {
	Id            int     `gorm:"primaryKey"`
	Name          string  `gorm:"index;not null"`
	Description   *string `gorm:"default:null"`
	CategoryId    int     `gorm:"index;not null"`
	SalesPrice    decimal.Decimal
	PurchasePrice decimal.Decimal
	UnitId        int `gorm:"index;not null"`
	Category      Category
	Unit          Unit
	HasShopId
	HasIsActive
}

type ItemResource struct {
	Id            int             `json:"id"`
	Name          string          `json:"name"`
	SalesPrice    decimal.Decimal `json:"sales_price"`
	PurchasePrice decimal.Decimal `json:"purchase_price"`
	CategoryName  string          `json:"category_name"`
	UnitName      string          `json:"unit_name"`
	UnitSymbol    string          `json:"unit_symbol"`
}

type ItemDetailResource struct {
	Id            int             `json:"id"`
	Name          string          `json:"name"`
	Description   *string         `json:"description"`
	SalesPrice    decimal.Decimal `json:"sales_price"`
	PurchasePrice decimal.Decimal `json:"purchase_price"`
	Category      struct {
		Id       int    `json:"id"`
		Name     string `json:"name"`
		IsActive bool   `json:"is_active"`
	} `json:"category"`
	Unit struct {
		Id       int    `json:"id"`
		Name     string `json:"name"`
		Symbol   string `json:"symbol"`
		IsActive bool   `json:"is_active"`
	} `json:"unit"`
	HasIsActive
	HasShopId
}

func FetchItemResources(db *gorm.DB, shopId string) ([]ItemResource, error) {
	var items []Item
	if err := db.Preload("Category").Preload("Unit").Where("shop_id = ? AND is_active = 1", shopId).Find(&items).Error; err != nil {
		return nil, err
	}
	results := make([]ItemResource, 0, len(items))
	for _, item := range items {
		results = append(results, ItemResource{
			Id:            item.Id,
			Name:          item.Name,
			SalesPrice:    item.SalesPrice,
			PurchasePrice: item.PurchasePrice,
			CategoryName:  item.Category.Name,
			UnitName:      item.Unit.Name,
			UnitSymbol:    item.Unit.Symbol,
		})
	}
	return results, nil
}

func FetchInactiveItemResources(db *gorm.DB, shopId string) ([]ItemResource, error) {
	var items []Item
	if err := db.Preload("Category").Preload("Unit").Where("shop_id = ? AND is_active = 0", shopId).Find(&items).Error; err != nil {
		return nil, err
	}
	results := make([]ItemResource, 0, len(items))
	for _, item := range items {
		results = append(results, ItemResource{
			Id:            item.Id,
			Name:          item.Name,
			SalesPrice:    item.SalesPrice,
			PurchasePrice: item.PurchasePrice,
			CategoryName:  item.Category.Name,
			UnitName:      item.Unit.Name,
			UnitSymbol:    item.Unit.Symbol,
		})
	}
	return results, nil
}

// func (readers *ReadServices) AfterUpdateItem(db *gorm.DB, shopId string, id int) error {
// 	if err := readers.ItemGetService.CleanCache(id); err != nil {
// 		return err
// 	}
// 	if err := readers.ItemListService.CleanCache(shopId); err != nil {
// 		return err
// 	}
// 	return nil
// }

type ItemSearch struct {
	Search           string
	CategoryId       int
	UnitId           int
	SalesPriceMin    *decimal.Decimal
	SalesPriceMax    *decimal.Decimal
	PurchasePriceMin *decimal.Decimal
	PurchasePriceMax *decimal.Decimal
}

func (s *ItemService) SearchItems(ctx context.Context, shopId string, search *ItemSearch) ([]ItemResource, error) {
	var items []Item
	dbCtx := s.db.WithContext(ctx).Where("shop_id = ?", shopId)
	if search.Search != "" {
		dbCtx.Where("name LIKE ? OR description LIKE ?",
			"%"+search.Search+"%",
			"%"+search.Search+"%",
		)
	}
	if search.CategoryId > 0 {
		dbCtx.Where("category_id = ?", search.CategoryId)
	}
	if search.UnitId > 0 {
		dbCtx.Where("unit_id = ?", search.UnitId)
	}
	if search.SalesPriceMin != nil {
		if search.SalesPriceMax == nil {
			dbCtx.Where("sales_price >= ?", search.SalesPriceMin)
		} else {
			dbCtx.Where("sales_price >= ? AND sales_price <= ?", search.SalesPriceMin, search.SalesPriceMax)
		}
	}
	if search.PurchasePriceMin != nil {
		if search.PurchasePriceMax == nil {
			dbCtx.Where("purchase_price >= ?", search.PurchasePriceMin)
		} else {
			dbCtx.Where("purchase_price >= ? AND purchase_price <= ?", search.PurchasePriceMin, search.PurchasePriceMax)
		}
	}

	if err := dbCtx.Preload("Category").Preload("Unit").Find(&items).Error; err != nil {
		return nil, err
	}

	results := make([]ItemResource, 0, len(items))
	for _, item := range items {
		results = append(results, ItemResource{
			Id:            item.Id,
			Name:          item.Name,
			SalesPrice:    item.SalesPrice,
			PurchasePrice: item.PurchasePrice,
			CategoryName:  item.Category.Name,
			UnitName:      item.Unit.Name,
			UnitSymbol:    item.Unit.Symbol,
		})
	}
	return results, nil
}

type ItemService struct {
	db     *gorm.DB
	getter services.Getter[ItemDetailResource]
	lister services.Lister[ItemResource]
}

func NewItemService(db *gorm.DB, cache services.CacheService) *ItemService {
	return &ItemService{
		db: db,
		getter: &customGetService[ItemDetailResource]{
			db:          db,
			cache:       cache,
			cachePrefix: "items",
			cacheLength: time.Hour * 127,
			fetch: func(db *gorm.DB, id int) (ItemDetailResource, error) {
				var item Item
				if err := db.Preload("Category").Preload("Unit").First(&item, id).Error; err != nil {
					return ItemDetailResource{}, err
				}
				result := ItemDetailResource{
					Id:            item.Id,
					Name:          item.Name,
					SalesPrice:    item.SalesPrice,
					PurchasePrice: item.PurchasePrice,
					Description:   item.Description,
				}
				result.ShopId = item.ShopId
				result.IsActive = item.IsActive
				result.Category.Id = item.Category.Id
				result.Category.Name = item.Category.Name
				result.Category.IsActive = item.Category.IsActive

				result.Unit.Id = item.Unit.Id
				result.Unit.Name = item.Unit.Name
				result.Unit.Symbol = item.Unit.Symbol
				result.Unit.IsActive = item.Unit.IsActive
				return result, nil
			},
		},
		lister: &customListService[ItemResource]{
			db:          db,
			cache:       cache,
			cachePrefix: "ItemList",
			cacheLength: forever,
			fetch:       FetchItemResources,
		},
	}
}

func (input *Item) validate(db *gorm.DB, shopId string, id int) error {
	shopFilter := NewShopFilter(shopId)
	if err := Validate(db,
		NewExistsRule("items", id, badRequest("item not found"), shopFilter).When(id > 0),
		NewUniqueRule("items", "name", input.Name, id, badRequest("duplicate item name"), shopFilter),
		NewExistsRule("units", input.UnitId, badRequest("unit not found"), shopFilter),
		NewExistsRule("categories", input.CategoryId, badRequest("category not found"), shopFilter),
	); err != nil {
		return err
	}
	return nil
}
func (s *ItemService) Store(ctx context.Context, input *Item, shopId string) (*Item, error) {
	if err := input.validate(s.db.WithContext(ctx), shopId, 0); err != nil {
		return nil, err
	}

	tx := s.db.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := tx.Create(&input).Error; err != nil {
		return nil, err
	}

	if err := s.lister.CleanCache(shopId); err != nil {
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	return input, nil
}

func (s *ItemService) cleanCacheAfterInsanceUpdate(ctx context.Context, id int, shopId string) error {
	if err := s.getter.CleanCache(id); err != nil {
		return err
	}
	if err := s.lister.CleanCache(shopId); err != nil {
		return err
	}
	return nil
}

func (s *ItemService) Update(ctx context.Context, input *Item, id int, shopId string) (*Item, error) {

	if err := input.validate(s.db.WithContext(ctx), shopId, id); err != nil {
		return nil, err
	}

	item, err := first[Item](s.db.WithContext(ctx), shopId, id)
	if err != nil {
		return nil, err
	}

	tx := s.db.WithContext(ctx).Begin()
	defer tx.Rollback()

	updates := map[string]any{
		"Name":          input.Name,
		"CategoryId":    input.CategoryId,
		"UnitId":        input.UnitId,
		"PurchasePrice": input.PurchasePrice,
		"SalesPrice":    input.SalesPrice,
	}
	if input.Description != nil {
		updates["Description"] = input.Description
	}
	if err := tx.Model(&item).Updates(updates).Error; err != nil {
		return nil, err
	}
	if err := s.cleanCacheAfterInsanceUpdate(ctx, id, shopId); err != nil {
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	return item, nil
}
func (s *ItemService) Delete(ctx context.Context, id int, shopId string) (*Item, error) {

	if err := Validate(s.db.WithContext(ctx),
		NewExistsRule("items", id, notFound("item not found"), NewShopFilter(shopId))); err != nil {
		return nil, err
	}

	item, err := first[Item](s.db.WithContext(ctx), shopId, id)
	if err != nil {
		return nil, err
	}

	tx := s.db.WithContext(ctx)
	defer tx.Rollback()

	if err := tx.Delete(&item).Error; err != nil {
		return nil, err
	}

	if err := s.cleanCacheAfterInsanceUpdate(ctx, id, shopId); err != nil {
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return item, nil
}

func (s *ItemService) Get(ctx context.Context, id int, shopId string) (*ItemDetailResource, bool, error) {

	return s.getter.Get(shopId, id)
}
func (s *ItemService) List(ctx context.Context, shopId string) ([]ItemResource, error) {
	return s.lister.List(shopId)
}
