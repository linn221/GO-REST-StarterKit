package models

import (
	"linn221/shop/services"

	"gorm.io/gorm"
)

type CrudServices struct {
	CategoryService *CategoryService
	ItemService     *ItemService
	UnitService     *UnitService
	UserService     *UserService
}

func NewServices(db *gorm.DB, cache services.CacheService) *CrudServices {
	itemService := NewItemService(db, cache)
	return &CrudServices{
		CategoryService: NewCategoryService(db, cache, itemService),
		UserService:     NewUserService(db),
		ItemService:     itemService,
		UnitService:     NewUnitService(db, cache, itemService),
	}
}
