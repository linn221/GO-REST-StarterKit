package models

type Category struct {
	Id          int     `gorm:"primaryKey"`
	Name        string  `gorm:"index;not null"`
	Description *string `gorm:"default:null"`
	HasShopId
	HasIsActive
}

// func (cat *Category) AfterCreate(db *gorm.DB) error {
// 	if err := ; err != nil {
// 		return err
// 	}
