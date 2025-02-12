package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"go.uber.org/zap"

	"shop/internal/adapter/password"
	repo "shop/internal/adapter/repo/mock"
	"shop/internal/app/entity"
	"shop/internal/app/port"
	"shop/internal/app/usecase"
)

func TestLogin(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	authRepo := repo.NewMockAuthRepo(ctrl)
	passHasher := password.NewMockPassHasher(ctrl)
	userUsecase := usecase.NewUser(authRepo, nil, nil, nil, passHasher, nil)

	login := entity.Login{
		Password: "test",
		Username: "test",
	}

	passHasher.EXPECT().Hash("test").Return("hash-test", nil).AnyTimes()
	passHasher.EXPECT().Compare("test", "hash-test").Return(true).AnyTimes()

	ctx := context.Background()
	authRepo.EXPECT().ReadPassword(ctx, "test").Return("hash-test", nil).AnyTimes()

	_, err := userUsecase.Auth(ctx, login)
	if err != nil {
		t.Error(err)
	}
}

func TestLoginWithInvalidPassword(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Error(err)
	}

	authRepo := repo.NewMockAuthRepo(ctrl)
	passHasher := password.NewHasherBcrypt(logger.Named("pass-hasher"))
	userUsecase := usecase.NewUser(authRepo, nil, nil, nil, passHasher, nil)

	login := entity.Login{
		Password: "invalid",
		Username: "test",
	}

	hash, err := passHasher.Hash("test")
	if err != nil {
		t.Error(err)
	}

	ctx := context.Background()

	authRepo.EXPECT().ReadPassword(ctx, "test").Return(hash, nil).AnyTimes()

	_, err = userUsecase.Auth(ctx, login)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidPassword) {
			return
		}

		t.Error(err)
	}
}

func TestRegister(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	authRepo := repo.NewMockAuthRepo(ctrl)
	passHasher := password.NewMockPassHasher(ctrl)
	userUsecase := usecase.NewUser(authRepo, nil, nil, nil, passHasher, nil)

	login := entity.Login{
		Password: "new",
		Username: "new",
	}

	passHasher.EXPECT().Hash("new").Return("hash-new", nil).AnyTimes()

	hash, err := passHasher.Hash(login.Password)
	if err != nil {
		t.Error(err)
	}

	loginWithHash := entity.Login{
		Username: "new",
		Password: hash,
	}

	ctx := context.Background()

	authRepo.EXPECT().ReadPassword(ctx, "new").Return("", port.ErrNotFound).AnyTimes()
	authRepo.EXPECT().CreateUser(ctx, loginWithHash).Return(nil).AnyTimes()

	_, err = userUsecase.Auth(ctx, login)
	if err != nil {
		t.Error(err)
	}

	// TODO: Check that JWT is created with right username
}
