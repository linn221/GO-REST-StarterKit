package services

import (
	"context"
	"linn221/shop/models"
	"linn221/shop/utils"

	"gorm.io/gorm"
)

type UserCruder interface {
	GetUser(context.Context, int) (*models.User, *ServiceError)
	ChangePassword(context.Context, int, string, string) *ServiceError
	UpdateInfo(context.Context, int, string, string, *string) *ServiceError
}

type userService struct {
	db    *gorm.DB
	cache CacheService
}

func (s *userService) GetUser(ctx context.Context, id int) (*models.User, *ServiceError) {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, systemErr(err)
	}

	return &user, nil
}

func (s *userService) ChangePassword(ctx context.Context, id int, oldPasword string, newPassword string) *ServiceError {
	user, errs := s.GetUser(ctx, id)
	if errs != nil {
		return errs
	}
	if err := utils.ComparePassword(user.Password, oldPasword); err != nil {
		return clientErr("passwords do not match")
	}
	hashed, err := utils.HashPassword(newPassword)
	if err != nil {
		return systemErr(err)
	}
	if err := s.db.Model(&user).UpdateColumn("password", hashed).Error; err != nil {
		return systemErr(err)
	}
	return nil
}

func (s *userService) UpdateInfo(ctx context.Context, id int, username string, email string, phoneNo *string) *ServiceError {
	user, errs := s.GetUser(ctx, id)
	if errs != nil {
		return errs
	}
	updates := map[string]any{
		"Username": username,
		"Email":    email,
	}
	if phoneNo != nil {
		updates["PhoneNo"] = *phoneNo
	}
	if err := Validate(s.db,
		NewUniqueRule("users", "username", username, id, "duplicate username"),
		NewUniqueRule("users", "email", email, id, "duplicate email"),
		NewUniqueRule("users", "phone_no", phoneNo, id, "duplicate phone number").When(phoneNo != nil),
	); err != nil {
		return err
	}
	if err := s.db.Model(&user).Updates(updates).Error; err != nil {
		return systemErr(err)
	}

	return nil
}
