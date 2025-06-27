package models

import (
	"context"
	"linn221/shop/services"

	"gorm.io/gorm"
)

type Unit struct {
	Id          int     `gorm:"primaryKey"`
	Name        string  `gorm:"index;not null"`
	Symbol      string  `gorm:"index;not null"`
	Description *string `gorm:"default:null"`
	HasShopId
	HasIsActive
}

type UnitResource struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

type UnitDetailResource struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Symbol      string  `json:"symbol"`
	Description *string `json:"description"`
	HasShopId
	HasIsActive
}

type UnitService struct {
	db                *gorm.DB
	getter            services.Getter[UnitDetailResource]
	lister            services.Lister[UnitResource]
	cleanRelatedCache func(ctx context.Context, id int, shopId string) error
}

func NewUnitService(db *gorm.DB, cache services.CacheService, itemService *ItemService) *UnitService {
	return &UnitService{
		db: db,
		getter: &defaultGetService[UnitDetailResource]{
			db:          db,
			cache:       cache,
			table:       "units",
			cachePrefix: "Unit",
			cacheLength: forever,
		},
		lister: &defaultListService[UnitResource]{
			db:          db,
			cache:       cache,
			table:       "units",
			cachePrefix: "UnitList",
			cacheLength: forever,
		},
		cleanRelatedCache: func(ctx context.Context, id int, shopId string) error {
			var relatedItemIds []int
			if err := db.WithContext(ctx).Model(&Item{}).Where("unit_id = ?", id).Pluck("id", &relatedItemIds).Error; err != nil {
				return err
			}
			for _, itemId := range relatedItemIds {
				if err := itemService.getter.CleanCache(itemId); err != nil {
					return err
				}
			}
			if err := itemService.lister.CleanCache(shopId); err != nil {
				return err
			}
			return nil
		},
	}
}

func (input *Unit) validate(db *gorm.DB, shopId string, id int) error {

	shopFilter := NewShopFilter(shopId)
	if err := Validate(db,
		NewExistsRule("units", id, notFound("unit not found"), shopFilter).When(id > 0),
		NewUniqueRule("units", "name", input.Name, id, badRequest("duplicate name"), shopFilter),
		NewUniqueRule("units", "symbol", input.Symbol, id, badRequest("duplicate symbol"), shopFilter),
	); err != nil {
		return err
	}
	return nil
}

func (s *UnitService) cleanCacheAfterInstanceUpdate(ctx context.Context, id int, shopId string) error {

	if err := s.getter.CleanCache(id); err != nil {
		return err
	}

	if err := s.lister.CleanCache(shopId); err != nil {
		return err
	}

	if err := s.cleanRelatedCache(ctx, id, shopId); err != nil {
		return err
	}

	return nil
}

func (s *UnitService) Store(ctx context.Context, input *Unit, shopId string) (*Unit, error) {
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

func (s *UnitService) Update(ctx context.Context, input *Unit, id int, shopId string) (*Unit, error) {
	if err := input.validate(s.db.WithContext(ctx), shopId, id); err != nil {
		return nil, err
	}

	unit, err := first[Unit](s.db.WithContext(ctx), shopId, id)
	if err != nil {
		return nil, err
	}

	tx := s.db.WithContext(ctx).Begin()
	defer tx.Rollback()

	updates := map[string]any{
		"Name":   input.Name,
		"Symbol": input.Symbol,
	}
	if input.Description != nil {
		updates["Description"] = input.Description
	}

	if err := tx.Model(&unit).Updates(updates).Error; err != nil {
		return nil, err
	}

	if err := s.cleanCacheAfterInstanceUpdate(ctx, id, shopId); err != nil {
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return input, nil
}

func (s *UnitService) Delete(ctx context.Context, id int, shopId string) (*Unit, error) {
	if err := Validate(s.db.WithContext(ctx),
		NewExistsRule("units", id, ErrNotFound, NewShopFilter(shopId)),
		NewNoResultRule("items", badRequest("unit is used in items"), NewFilter("unit_id = ?", id)),
	); err != nil {
		return nil, err
	}

	unit, err := first[Unit](s.db.WithContext(ctx), shopId, id)
	if err != nil {
		return nil, err
	}

	tx := s.db.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := tx.Delete(&unit).Error; err != nil {
		return nil, err
	}

	if err := s.cleanCacheAfterInstanceUpdate(ctx, id, shopId); err != nil {
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return unit, nil
}

func (s *UnitService) Get(ctx context.Context, id int, shopId string) (*UnitDetailResource, bool, error) {
	return s.getter.Get(shopId, id)
}

func (s *UnitService) List(ctx context.Context, shopId string) ([]UnitResource, error) {
	return s.lister.List(shopId)
}
