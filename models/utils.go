package models

import (
	"gorm.io/gorm"
)

func first[T any](db *gorm.DB, shopId string, id int) (*T, error) {
	var v T
	err := db.Where("shop_id = ?", shopId).First(&v, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &v, nil
}
