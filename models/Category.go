package models

import (
	"context"
	"linn221/shop/services"

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

func (input *Category) validate(db *gorm.DB, shopId string, id int) error {
	shopFilter := NewShopFilter(shopId)
	return Validate(db,
		NewExistsRule("categories", id, notFound("category not found"), shopFilter).When(id > 0),
		NewUniqueRule("categories", "name", input.Name, id, badRequest("duplicate category name"), NewShopFilter(shopId)),
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

func (s *CategoryService) Store(ctx context.Context, input *Category, shopId string) (*Category, error) {
	err := input.validate(s.db.WithContext(ctx), shopId, 0)
	if err != nil {
		return nil, err
	}
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&input).Error; err != nil {
			return err
		}
		if err := s.lister.CleanCache(shopId); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return input, nil
}

func (s *CategoryService) Update(ctx context.Context, input *Category, id int, shopId string) (*Category, error) {

	err := input.validate(s.db.WithContext(ctx), shopId, 0)
	if err != nil {
		return nil, err
	}

	category, err := first[Category](s.db.WithContext(ctx), shopId, id)
	if err != nil {
		return nil, err
	}

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		updates := map[string]any{
			"Name":        input.Name,
			"Description": input.Description,
		}
		if err := tx.Model(&category).Updates(updates).Error; err != nil {
			return err
		}
		if err := s.getter.CleanCache(id); err != nil {
			return err
		}
		if err := s.lister.CleanCache(shopId); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return category, nil
}

func (s *CategoryService) Delete(ctx context.Context, id int, shopId string) (*Category, error) {
	category, err := first[Category](s.db.WithContext(ctx), shopId, id)
	if err != nil {
		return nil, err
	}

	if err := Validate(s.db.WithContext(ctx),
		NewNoResultRule("items", badRequest("category has been used in items"), NewFilter("category_id = ?", id)),
	); err != nil {
		return nil, err
	}

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&category).Error; err != nil {
			return err
		}
		if err := s.getter.CleanCache(id); err != nil {
			return err
		}
		if err := s.lister.CleanCache(shopId); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return category, nil
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
// func (readers *ReadServices) AfterCategoryUpdate(db *gorm.DB, shopId string, id int) error {

// 	if err := readers.CategoryGetService.CleanCache(id); err != nil {
// 		return err
// 	}
// 	if err := readers.CategoryListService.CleanCache(shopId); err != nil {
// 		return err
// 	}
// 	var relatedItemIds []int
// 	if err := db.Model(&Item{}).Where("category_id = ?", id).Pluck("id", &relatedItemIds).Error; err != nil {
// 		return err
// 	}
// 	for _, itemId := range relatedItemIds {
// 		if err := readers.ItemGetService.CleanCache(itemId); err != nil {
// 			return err
// 		}
// 	}
// 	if err := readers.ItemListService.CleanCache(shopId); err != nil {
// 		return err
// 	}
// 	return nil

// }
