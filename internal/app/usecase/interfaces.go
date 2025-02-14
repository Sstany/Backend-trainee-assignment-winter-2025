package usecase

import (
	"context"

	"shop/internal/app/entity"
)

type UserUseCase interface {
	Auth(ctx context.Context, login entity.Login) (*entity.Token, error)
	Send(ctx context.Context, reqSend entity.SendCoinRequest) error
	Info(ctx context.Context, request entity.InfoRequest) (*entity.InfoResponse, error)
	Buy(ctx context.Context, item entity.ItemRequest) error
}
