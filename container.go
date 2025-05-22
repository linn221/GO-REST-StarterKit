package main

import (
	"linn221/shop/config"
	"linn221/shop/models"

	"gorm.io/gorm"
)

type Container struct {
	DB             *gorm.DB
	Cache          *config.RedisCache
	ImageDirectory string
	Readers        *models.ReadServices
}
