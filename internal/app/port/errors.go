package port

import "errors"

var (
	ErrAlreadyRegistred = errors.New("user alredy exists")

	ErrNotFound = errors.New("user login not found")

	ErrInsufficientBalance = errors.New("not enought coins")

	ErrReveicerNotExists = errors.New("wrong receiver name")
)
