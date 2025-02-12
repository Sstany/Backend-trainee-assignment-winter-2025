package port

import (
	"context"
	"errors"

	"shop/internal/app/entity"
)

var (
	ErrNotFound = errors.New("user login not found")
)

//go:generate mockgen -destination ../../adapter/repo/auth_mock.go -package repo -source ./auth.go

type AuthRepo interface {
	ReadPassword(ctx context.Context, username string) (string, error)
	CreateUser(ctx context.Context, login entity.Login) error
}
