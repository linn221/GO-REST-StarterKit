package main

import (
	"linn221/shop/models"
	"linn221/shop/services"

	"gorm.io/gorm"
)

type Container struct {
	DB             *gorm.DB
	Cache          services.CacheService
	ImageDirectory string
	Readers        *models.ReadServices
}
