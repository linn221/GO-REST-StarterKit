package services

import (
	"github.com/go-playground/validator"
	"gorm.io/gorm"
)

type Container struct {
	DB       *gorm.DB
	Cache    CacheService
	Validate *validator.Validate
}
