package models

import (
	"linn221/shop/services"
	"net/http"

	"gorm.io/gorm"
)

func first[T any](db *gorm.DB, shopId string, id int) (*T, *services.ServiceError) {
	var v T
	err := db.Where("shop_id = ?", shopId).First(&v, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &services.ServiceError{Code: http.StatusNotFound, Err: err}
		}
		return nil, services.SystemErr(err)
	}

	return &v, nil
}
