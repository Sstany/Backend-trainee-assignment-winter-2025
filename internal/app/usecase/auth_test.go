package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"

	"shop/internal/adapter/password"
	repo "shop/internal/adapter/repo/mock"
	"shop/internal/app/entity"
	"shop/internal/app/usecase"
)

func TestLogin(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	authRepo := repo.NewMockAuthRepo(ctrl)
	passHasher := password.NewMockPassHasher(ctrl)
	secretRepo := repo.NewMockSecretRepo(ctrl)
	balanceRepo := repo.NewMockUserBalanceRepo(ctrl)

	authUsecase, _ := usecase.NewAuth(authRepo, balanceRepo, passHasher, secretRepo, nil)

	login := entity.Login{
		Password: "testtest",
		Username: "test",
	}
	ctx := context.Background()

	authRepo.EXPECT().ReadPassword(ctx, "test").Return("hash-test", nil).AnyTimes()
	authRepo.EXPECT().CreateUser(ctx, login).Return(nil).AnyTimes()
	passHasher.EXPECT().Hash("testtest").Return("hash-test", nil).AnyTimes()
	balanceRepo.EXPECT().Create(ctx, "test", 1000).Return(nil).AnyTimes()

	passHasher.EXPECT().Compare("testtest", "hash-test").Return(true).AnyTimes()
	var token entity.Token
	token = "str"
	secretRepo.EXPECT().CreateToken("test").Return(&token, nil).AnyTimes()

	_, err := authUsecase.Auth(ctx, login)
	if err != nil {
		t.Error(err)
	}
}

func TestLoginWithInvalidPassword(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	authRepo := repo.NewMockAuthRepo(ctrl)
	passHasher := password.NewMockPassHasher(ctrl)
	secretRepo := repo.NewMockSecretRepo(ctrl)
	balanceRepo := repo.NewMockUserBalanceRepo(ctrl)

	authUsecase, _ := usecase.NewAuth(authRepo, balanceRepo, passHasher, secretRepo, nil)

	login := entity.Login{
		Password: "testtest",
		Username: "test",
	}
	ctx := context.Background()

	authRepo.EXPECT().ReadPassword(ctx, "test").Return("", usecase.ErrInvalidPassword)
	passHasher.EXPECT().Hash("testtest").Return("", nil).AnyTimes()

	_, err := authUsecase.Auth(ctx, login)
	if !errors.Is(err, usecase.ErrInvalidPassword) {
		t.Error(err)
	}
}
