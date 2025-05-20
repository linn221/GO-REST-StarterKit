package services

import (
	"linn221/shop/models"

	"gorm.io/gorm"
)

type CategoryCruder interface {
	CreateCategory(input *models.Category, db *gorm.DB, cache CacheService) (*models.Category, *ServiceError)
}
