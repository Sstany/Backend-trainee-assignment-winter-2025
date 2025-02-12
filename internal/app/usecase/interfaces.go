package usecase

import (
	"context"

	"shop/internal/app/entity"
)

type UserUseCase interface {
	Auth(ctx context.Context, login entity.Login) (*entity.Token, error)
	// Send() error
	// Info() (*entity.UserInfo, error)
	// Buy(item *entity.Product) error
}
