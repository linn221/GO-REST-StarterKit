package services

import (
	"encoding/json"
	"errors"
	"linn221/shop/models"
	"linn221/shop/utils"
	"net/http"

	"gorm.io/gorm"
)

type userService struct{}
type ServiceError struct {
	Err  error
	Code int
}

func (se *ServiceError) Respond(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(se.Code)
	return json.NewEncoder(w).Encode(map[string]any{
		"error": se.Err.Error(),
	})
}

func systemErr(err error) *ServiceError {
	return &ServiceError{
		Err:  err,
		Code: http.StatusInternalServerError,
	}
}
func systemErrString(s string) *ServiceError {
	return systemErr(errors.New(s))
}
func clientErr(s string) *ServiceError {
	return &ServiceError{
		Err:  errors.New(s),
		Code: http.StatusBadRequest,
	}
}

func (s *userService) GetUser(id int, db *gorm.DB, cache CacheService) (*models.User, *ServiceError) {
	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		return nil, systemErr(err)
	}

	return &user, nil
}

func (s *userService) ChangePassword(id int, oldPasword string, newPassword string, db *gorm.DB, cache CacheService) *ServiceError {
	user, errs := s.GetUser(id, db, cache)
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
	if err := db.Model(&user).UpdateColumn("password", hashed).Error; err != nil {
		return systemErr(err)
	}
	return nil
}

func (s *userService) UpdateInfo(id int, name string, username string, email string, phoneNo *string, db *gorm.DB, cache CacheService) *ServiceError {
	user, errs := s.GetUser(id, db, cache)
	if errs != nil {
		return errs
	}
	updates := map[string]any{
		"Name":     name,
		"Username": username,
		"Email":    email,
	}
	if phoneNo != nil {
		updates["PhoneNo"] = *phoneNo
	}
	if err := db.Model(&user).Updates(updates).Error; err != nil {
		return systemErr(err)
	}

	return nil
}
