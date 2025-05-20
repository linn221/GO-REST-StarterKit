package services

import (
	"context"
	"linn221/shop/models"
)

type CategoryCruder interface {
	CreateCategory(ctx context.Context, input *models.Category) (*models.Category, *ServiceError)
	UpdateCategory(ctx context.Context, id int, input *models.Category, existing *models.Category) (*models.Category, *ServiceError)
	DeleteCategory(ctx context.Context, id int, existing *models.Category) (*models.Category, *ServiceError)
	GetCategory(ctx context.Context, id int) (*models.Category, *ServiceError)
	ListCategories(ctx context.Context) ([]*models.Category, *ServiceError)
}
