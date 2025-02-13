package usecase

import (
	"context"

	"shop/internal/app/entity"
)

type UserUseCase interface {
	Auth(ctx context.Context, login entity.Login) (*entity.Token, error)
	Send(ctx context.Context, reqSend entity.SendCoinRequest) error
	//Info() (*entity.UserInfo, error)
	Buy(ctx context.Context, item entity.ItemRequest) error
}
