package models

type Category struct {
	Id   int    `gorm:"primaryKey"`
	Name string `gorm:"index;not null"`
	HasShopId
	HasIsActive
}
