package models

import (
	"context"
	"sync"

	"gorm.io/gorm"
)

type User struct {
	Id       int    `gorm:"primaryKey"`
	Username string `gorm:"unique;not null"`
	Email    string `gorm:"unique"`
	PhoneNo  string `gorm:"unique"`
	Password string `gorm:"index;not null"`
	HasIsActive
	HasShopId
}

type UserService struct {
	mu sync.Mutex
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		db: db,
	}
}

func (u *UserService) UpdateInfo(ctx context.Context, shopId string, userId int, input *User) (*User, error) {

	u.mu.Lock()
	defer u.mu.Unlock()

	var user User
	if err := u.db.WithContext(ctx).First(&user, userId).Error; err != nil {
		return nil, err
	}

	updates := map[string]any{
		"Username": input.Username,
		"Email":    input.Email,
	}
	if input.PhoneNo != "" {
		updates["PhoneNo"] = input.PhoneNo
	}

	shopFilter := NewShopFilter(shopId)
	if err := Validate(u.db.WithContext(ctx),
		NewUniqueRule("users", "username", input.Username, userId, badRequest("duplicate username"), shopFilter),
		NewUniqueRule("users", "email", input.Email, userId, badRequest("duplicate email"), shopFilter),
		NewUniqueRule("users", "phone_no", input.PhoneNo, userId, badRequest("duplicate phone number"), shopFilter).When(input.PhoneNo != ""),
	); err != nil {
		return nil, err
	}

	if err := u.db.WithContext(ctx).Model(&user).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
