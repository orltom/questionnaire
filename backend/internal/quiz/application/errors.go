package application

import "errors"

var (
	ErrInvalidArguments = errors.New("invalid arguments")
	ErrPersistence      = errors.New("persistence error")
	ErrEntityNotFound   = errors.New("entity not found")
	ErrForbidden        = errors.New("forbidden")
)
