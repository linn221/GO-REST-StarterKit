package models

type Category struct {
	Id          int     `gorm:"primaryKey"`
	Name        string  `gorm:"index;not null"`
	Description *string `gorm:"default:null"`
	HasShopId
	HasIsActive
}

type CategoryResource struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	HasShopId
}

type CategoryDetailResource struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	HasIsActive
	HasShopId
}

// func (cat *Category) AfterCreate(db *gorm.DB) error {
// 	if err := ; err != nil {
// 		return err
// 	}
