package models

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound   = errors.New("record not found")
	ErrBadRequest = errors.New("bad request")
	ErrDb         = errors.New("db error")
)

func badRequest(s string) error {
	return fmt.Errorf("%w: %s", ErrBadRequest, s)
}

func wrapError(wrapedErr, err error) error {
	return fmt.Errorf("%w: %v", wrapedErr, err)
}

func notFound(s string) error {
	if s == "" {
		return ErrNotFound
	}
	return fmt.Errorf("%w: %s", ErrNotFound, s)
}

func dbError(err error) error {
	return wrapError(ErrDb, err)
}
