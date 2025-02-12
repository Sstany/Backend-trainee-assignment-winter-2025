package port

import "errors"

var (
	ErrNotFound = errors.New("user login not found")

	ErrInsufficientBalance = errors.New("not enought coins")
)
