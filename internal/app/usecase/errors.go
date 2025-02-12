package usecase

import "errors"

var (
	ErrInvalidPassword = errors.New("invalid password")

	ErrUnsafePassword = errors.New("password too short")
	ErrLongPassword   = errors.New("password too long")

	ErrItemNotExists = errors.New("item not exists")

	ErrInsufficientBalance = errors.New("not enought coins")
)
