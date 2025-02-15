package usecase

import "errors"

var (
	ErrInvalidPassword    = errors.New("invalid password")
	ErrUnsafePassword     = errors.New("password too short")
	ErrLongPassword       = errors.New("password too long")
	ErrInvalidSigningAlgo = errors.New("invalid signing algo")
	ErrTokenExpired       = errors.New("token expired")
	ErrInvalidToken       = errors.New("invalid token")

	ErrItemNotExists = errors.New("item not exists")

	ErrInsufficientBalance = errors.New("not enought coins")
	ErrWrongCoinAmount     = errors.New("coin amount <0")
)
