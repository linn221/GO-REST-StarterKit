package models

import (
	"linn221/shop/services"

	"gorm.io/gorm"
)

type ReadServices struct {
	CategoryGetService  services.Getter[Category]
	CategoryListService services.Lister[Category]
}

func NewReadServices(db *gorm.DB, cache services.CacheService) *ReadServices {
	return &ReadServices{
		CategoryGetService: &generalGetService[Category]{
			db:          db,
			cache:       cache,
			table:       "categories",
			cachePrefix: "Category",
			cacheLength: forever,
		},
		CategoryListService: &generalListService[Category]{
			db:          db,
			cache:       cache,
			table:       "categories",
			cachePrefix: "CategoryList",
			cacheLength: forever,
		},
	}
}
