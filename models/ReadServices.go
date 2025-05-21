package models

import (
	"linn221/shop/services"

	"gorm.io/gorm"
)

type ReadServices struct {
	CategoryGetService  services.Getter[CategoryResource]
	CategoryListService services.Lister[CategoryResource]
}

func NewReaders(db *gorm.DB, cache services.CacheService) *ReadServices {
	return &ReadServices{
		CategoryGetService: &defaultGetService[CategoryResource]{
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
	}
}
