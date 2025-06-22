package models

import (
	"context"
	"linn221/shop/services"
	"linn221/shop/validate"

	"gorm.io/gorm"
)

type Category struct {
	Id          int     `gorm:"primaryKey"`
	Name        string  `gorm:"index;not null"`
	Description *string `gorm:"default:null"`
	HasShopId
	HasIsActive
}

type CategoryResource struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	HasShopId
}

type CategoryDetailResource struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	HasIsActive
	HasShopId
}

type CategoryService struct {
	db     *gorm.DB
	getter services.Getter[CategoryDetailResource]
	lister services.Lister[CategoryResource]
}

func (input *Category) validate(db *gorm.DB, shopId string, id int) *services.ServiceError {
	shopFilter := validate.NewShopFilter(shopId)
	return validate.Validate(db,
		validate.NewExistsRule("categories", id, "category not found", shopFilter).When(id > 0),
		validate.NewUniqueRule("categories", "name", input.Name, id, "duplicate category name", validate.NewShopFilter(shopId)),
	)
}

func NewCategoryService(db *gorm.DB, cache services.CacheService) *CategoryService {
	return &CategoryService{
		db: db,
		getter: &defaultGetService[CategoryDetailResource]{
			db:          db,
			cache:       cache,
			table:       "categories",
			cachePrefix: "Category",
			cacheLength: forever,
		},
		lister: &defaultListService[CategoryResource]{
			db:          db,
			cache:       cache,
			table:       "categories",
			cachePrefix: "CategoryList",
			cacheLength: forever,
		},
	}
}

func (s *CategoryService) Store(ctx context.Context, input *Category, shopId string) (*Category, *services.ServiceError) {
	e := input.validate(s.db.WithContext(ctx), shopId, 0)
	if e != nil {
		return nil, e
	}
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&input).Error; err != nil {
			return err
		}
		if err := s.lister.CleanCache(shopId); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, services.SystemErr(err)
	}

	return input, nil
}

func (s *CategoryService) Update(ctx context.Context, input *Category, id int, shopId string) (*Category, *services.ServiceError) {

	e := input.validate(s.db.WithContext(ctx), shopId, 0)
	if e != nil {
		return nil, e
	}

	category, e := first[Category](s.db.WithContext(ctx), shopId, id)
	if e != nil {
		return nil, e
	}

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		return nil
	})

	if err != nil {
		return nil, services.SystemErr(err)
	}
}

func (s *CategoryService) Delete(ctx context.Context, id int, shopId int) (*Category, *services.ServiceError) {
	panic("2d")

}

func (s *CategoryService) Get(ctx context.Context, id int, shopId string) (*CategoryDetailResource, bool, error) {
	return s.getter.Get(shopId, id)
}

func (s *CategoryService) List(ctx context.Context, shopId string) ([]CategoryResource, error) {
	return s.lister.List(shopId)
}

// func (cat *Category) AfterCreate(db *gorm.DB) error {
// 	if err := ; err != nil {
// 		return err
// 	}

// clean related cache when category is updated
func (readers *ReadServices) AfterCategoryUpdate(db *gorm.DB, shopId string, id int) error {

	if err := readers.CategoryGetService.CleanCache(id); err != nil {
		return err
	}
	if err := readers.CategoryListService.CleanCache(shopId); err != nil {
		return err
	}
	var relatedItemIds []int
	if err := db.Model(&Item{}).Where("category_id = ?", id).Pluck("id", &relatedItemIds).Error; err != nil {
		return err
	}
	for _, itemId := range relatedItemIds {
		if err := readers.ItemGetService.CleanCache(itemId); err != nil {
			return err
		}
	}
	if err := readers.ItemListService.CleanCache(shopId); err != nil {
		return err
	}
	return nil

}
