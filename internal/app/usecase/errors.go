package usecase

import "errors"

var (
	ErrInvalidPassword = errors.New("invalid password")

	ErrItemNotExists = errors.New("item not exists")

	ErrInsufficientBalance = errors.New("not enought coins")
)
