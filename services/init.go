package services

import (
	"linn221/shop/config"

	"gorm.io/gorm"
)

func init() {
}

type Instance struct {
	UserService    *userService
	ImageService   ImageUploader
	CategoryCruder CategoryCruder
}

func NewServices(db *gorm.DB, cache CacheService) *Instance {
	dir := config.GetImageDirectory()
	return &Instance{
		UserService:  &userService{db: db, cache: cache},
		ImageService: &ImageUploadService{dir: dir, maxMemoryMB: 10, db: db},
	}
}
