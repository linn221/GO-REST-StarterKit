package models

type HasIsActive struct {
	IsActive bool `gorm:"default:true" json:"is_active"`
}

func (h HasIsActive) GetIsActive() bool {
	return h.IsActive
}

func (h *HasIsActive) SetActive() {
	h.IsActive = true
}

type HasUserId struct {
	UserId int `gorm:"index;not null" json:"user_id"`
}

func (h *HasUserId) GetUserId() int {
	return h.UserId
}

func (h *HasUserId) InjectUserId(userId int) {
	h.UserId = userId
}

type HasShopId struct {
	ShopId int `gorm:"index;not null"`
}

func (h *HasShopId) GetShopId() int {
	return h.ShopId
}

func (h *HasShopId) InjectShopId(userId int) {
	h.ShopId = userId
}
