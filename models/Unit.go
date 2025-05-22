package models

type Unit struct {
	Id          int     `gorm:"primaryKey"`
	Name        string  `gorm:"index;not null"`
	Symbol      string  `gorm:"index;not null"`
	Description *string `gorm:"default:null"`
	HasShopId
	HasIsActive
}

type UnitResource struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

type UnitDetailResource struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Symbol      string  `json:"symbol"`
	Description *string `json:"description"`
	HasIsActive
}
