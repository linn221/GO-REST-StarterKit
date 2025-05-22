package models

import (
	"linn221/shop/services"

	"gorm.io/gorm"
)

type ReadServices struct {
	CategoryGetService  services.Getter[CategoryDetailResource]
	CategoryListService services.Lister[CategoryResource]
	UnitGetService      services.Getter[UnitDetailResource]
	UnitListService     services.Lister[UnitResource]
}

func NewReaders(db *gorm.DB, cache services.CacheService) *ReadServices {
	return &ReadServices{
		CategoryGetService: &defaultGetService[CategoryDetailResource]{
			db:          db,
			cache:       cache,
			table:       "categories",
			cachePrefix: "Category",
			cacheLength: forever,
		},
		CategoryListService: &defaultListService[CategoryResource]{
			db:          db,
			cache:       cache,
			table:       "categories",
			cachePrefix: "CategoryList",
			cacheLength: forever,
		},
		UnitGetService: &defaultGetService[UnitDetailResource]{
			db:          db,
			cache:       cache,
			table:       "units",
			cachePrefix: "Unit",
			cacheLength: forever,
		},
		UnitListService: &defaultListService[UnitResource]{
			db:          db,
			cache:       cache,
			table:       "units",
			cachePrefix: "UnitList",
			cacheLength: forever,
		},
	}
}
