package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang-jwt/jwt"
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

	token := entity.Token("str")
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

func TestAuthenicateWithAccessTokenExpired(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	authRepo := repo.NewMockAuthRepo(ctrl)
	passHasher := password.NewMockPassHasher(ctrl)
	secretRepo := repo.NewMockSecretRepo(ctrl)
	balanceRepo := repo.NewMockUserBalanceRepo(ctrl)

	authUsecase, _ := usecase.NewAuth(authRepo, balanceRepo, passHasher, secretRepo, nil)

	//nolint:lll
	token := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mzk3NDk2NzcsImlzcyI6IiIsImp0aSI6ImVhMzIyNjlhLTliMGItNDEzYS1hOGMyLTEyZDY1MGQ1NWRhZCIsInVpZCI6IlRhbmkzIn0.wkbeIUAhcRUD1mi2kqA6spQobCxsGzINizjd0sWkIg3BgUwbvlERlJzNwpzvMy5dEn8sMAhwAl15qF17MX3JpAjzSh4uxTdKfXIGFjr1wvPLFu50Uc8lp4bivcALSf49E_N7rfsSKrEGEi7MNdFpRtjKQmgoGPf2lwEgqBFo3Gwbe2eJ_fP2yyNBZzDz9LoVOjjkC0944U9d2e7St2En-giSP1w7tUGWBSIgUzi6K17VV8UidD3UO0UScLOwMFRmYpG8AUbbTy5ymLUKxpzsHGN1yHfAB_8UcUj55B_TE7nJz1-Mgp_5KvMxRJPmxJSnV4fl76nWcsvJTXkJ5xNAaw"

	secretRepo.EXPECT().ParseJWT(token).Return(nil, usecase.ErrTokenExpired)

	_, _, err := authUsecase.AuthenticateWithAccessToken(token)
	if !errors.Is(err, usecase.ErrTokenExpired) {
		t.Error(err)
	}
}

func TestAuthenicateWithAccessToken(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	authRepo := repo.NewMockAuthRepo(ctrl)
	passHasher := password.NewMockPassHasher(ctrl)
	secretRepo := repo.NewMockSecretRepo(ctrl)
	balanceRepo := repo.NewMockUserBalanceRepo(ctrl)

	authUsecase, _ := usecase.NewAuth(authRepo, balanceRepo, passHasher, secretRepo, nil)

	//nolint:lll
	token := entity.Token("eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDgzMzY3NDUsImlzcyI6ImF2aXRvLXNob3AiLCJqdGkiOiI2MzA5ZmEzZi01N2FlLTQxMmUtYmYwZC03YmRmMTgwNGVhZTAiLCJ1aWQiOiJUYW5pMyJ9.qFzuGvXd5t_NVslwk4zHiG4hgzxz5QJ62UMxijkB5p5Rnf1oU_zmq1S0n2d7xfDkS9bBz_ptqaQmmkCsQVcBm8-NyAMJTz8uIhly0fj29nhZo_RnSvOdK_hqlRhIqr1pVKPGr5wAK_A5PgknMIVlSPc6j_xOlCydHZ_1RxzL4utidjoieYuo6oRyJNqr6F6xuK82I8K6xIA1M12uVmlN_43Zb-HWCX9--21f4ZhvjJrklutawh77GRJcxfUkkOPZJyblEdaoxYenTwd8ffyC70oXZvtfT6YPpxS7xXjdVfSkE3X6sz-7Wz_5-LlMAlrTcM2biHaBlJLIVWDwAn_Gbg")

	mc := jwt.MapClaims{
		"exp": float64(1748336745),
		"iss": "avito-shop",
		"jti": "6309fa3f-57ae-412e-bf0d-7bdf1804eae0",
		"uid": "Tani3",
	}

	secretRepo.EXPECT().ParseJWT(string(token)).Return(mc, nil)

	_, _, err := authUsecase.AuthenticateWithAccessToken(string(token))
	if err != nil {
		t.Error(err)
	}
}
