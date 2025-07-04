package models

import (
	"context"
	"fmt"
	"linn221/shop/utils"
	"math/rand"
	"sync"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	Id       int    `gorm:"primaryKey"`
	Username string `gorm:"unique;not null"`
	Email    string `gorm:"unique"`
	PhoneNo  string `gorm:"unique"`
	Password string `gorm:"index;not null" json:"-"`
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

func (u *UserService) Login(ctx context.Context, username string, password string) (*User, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	var user User
	if err := u.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, badRequest("invalid username or password")
		}
		return nil, err
	}

	if err := utils.ComparePassword(user.Password, password); err != nil {
		return nil, badRequest("invalid username or password")
	}
	return &user, nil
}

func (s *UserService) ChangePassword(ctx context.Context, userId int, oldPassword string, newPassword string) (*User, error) {

	s.mu.Lock()
	defer s.mu.Unlock()

	var user User
	if err := s.db.WithContext(ctx).First(&user, userId).Error; err != nil {
		return nil, err
	}

	if err := utils.ComparePassword(user.Password, oldPassword); err != nil {
		return nil, badRequest("password do not match")
	}
	hashed, err := utils.HashPassword(newPassword)
	if err != nil {
		return nil, err
	}
	if err := s.db.Model(&user).UpdateColumn("password", hashed).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserService) Register(ctx context.Context, name, email, password, phoneNo string) (*User, error) {

	u.mu.Lock()
	defer u.mu.Unlock()

	shopId := uuid.NewString()
	shop := Shop{
		Id:      shopId,
		Name:    name,
		LogoUrl: "",
		Email:   email,
		PhoneNo: phoneNo,
	}

	if err := Validate(u.db.WithContext(ctx),
		NewUniqueRule("users", "email", email, 0, badRequest("duplicate email"), nil),
		NewUniqueRule("users", "phone_no", phoneNo, 0, badRequest("duplicate phone number"), nil),
	); err != nil {
		return nil, err
	}

	tx := u.db.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := tx.Create(&shop).Error; err != nil {
		return nil, err
	}

	i := rand.Intn(100000)
	username := fmt.Sprintf("owner%d", i)
	password2, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := User{
		Username: username,
		Email:    email,
		PhoneNo:  phoneNo,
		Password: string(password2),
	}
	user.ShopId = shopId
	if err := tx.Create(&user).Error; err != nil {
		return nil, err
	}
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	return &user, nil
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

	if err := Validate(u.db.WithContext(ctx),
		NewUniqueRule("users", "username", input.Username, userId, badRequest("duplicate username"), nil),
		NewUniqueRule("users", "email", input.Email, userId, badRequest("duplicate email"), nil),
		NewUniqueRule("users", "phone_no", input.PhoneNo, userId, badRequest("duplicate phone number"), nil).When(input.PhoneNo != ""),
	); err != nil {
		return nil, err
	}

	if err := u.db.WithContext(ctx).Model(&user).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
