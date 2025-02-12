package port

import (
	"context"

	"shop/internal/app/entity"
)

//go:generate mockgen -destination ../../adapter/repo/mock/auth_mock.go -package repo -source ./auth.go

type AuthRepo interface {
	ReadPassword(ctx context.Context, username string) (string, error)
	CreateUser(ctx context.Context, login entity.Login) error
}
