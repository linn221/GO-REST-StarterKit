package services

import "gorm.io/gorm"

type Container struct {
	DB    *gorm.DB
	Cache CacheService
}
