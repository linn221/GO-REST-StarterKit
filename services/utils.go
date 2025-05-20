package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

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
func systemErrString(s string, args ...error) *ServiceError {
	if len(args) > 0 {
		err := args[0]
		return systemErr(errors.New(s + ": " + err.Error()))
	}
	return systemErr(errors.New(s))
}
func clientErr(s string) *ServiceError {
	return &ServiceError{
		Err:  errors.New(s),
		Code: http.StatusBadRequest,
	}
}

func NewSession(userId int, shopId string, cache CacheService) (string, error) {
	token := uuid.NewString()
	cacheVal := fmt.Sprintf("%d:%s", userId, shopId)
	cacheKey := "Token:" + token
	err := cache.SetValue(cacheKey, cacheVal, time.Hour*127)
	return token, err
}
