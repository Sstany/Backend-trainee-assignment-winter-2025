package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"

	repo "shop/internal/adapter/repo/mock"
	"shop/internal/app/entity"
	"shop/internal/app/port"
	"shop/internal/app/usecase"
)

func TestSend(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	shopRepo := repo.NewMockShopRepo(ctrl)
	balanceRepo := repo.NewMockUserBalanceRepo(ctrl)
	inventoryRepo := repo.NewMockUserInventoryRepo(ctrl)
	transactionController := repo.NewMockTransactionController(ctrl)
	userTransactionRepo := repo.NewMockUserTransactionRepo(ctrl)

	userUsecase := usecase.NewUser(shopRepo, balanceRepo, inventoryRepo, transactionController, userTransactionRepo)

	sendRequest := entity.SendCoinRequest{
		Amount:   200,
		FromUser: "test",
		ToUser:   "test2",
	}

	ctx := context.Background()
	tx := repo.NewMockTransaction(ctrl)
	transactionController.EXPECT().BeginTx(ctx).Return(tx, nil)

	balanceRepo.EXPECT().ChangeUserBalance(tx, -200, "test").Return(nil)
	balanceRepo.EXPECT().ChangeUserBalance(tx, 200, "test2").Return(nil)
	userTransactionRepo.EXPECT().SetUserTransaction(tx, sendRequest).Return(nil)

	tx.EXPECT().Commit().Return(nil)

	err := userUsecase.Send(ctx, sendRequest)
	if err != nil {
		t.Error(err)
	}
}

func TestSendWithInsufficientBalance(t *testing.T) {
	ctrl := gomock.NewController(t)

	shopRepo := repo.NewMockShopRepo(ctrl)
	balanceRepo := repo.NewMockUserBalanceRepo(ctrl)
	inventoryRepo := repo.NewMockUserInventoryRepo(ctrl)
	transactionController := repo.NewMockTransactionController(ctrl)
	userTransactionRepo := repo.NewMockUserTransactionRepo(ctrl)

	userUsecase := usecase.NewUser(shopRepo, balanceRepo, inventoryRepo, transactionController, userTransactionRepo)

	sendRequest := entity.SendCoinRequest{
		Amount:   200,
		FromUser: "test",
		ToUser:   "test2",
	}

	ctx := context.Background()
	tx := repo.NewMockTransaction(ctrl)
	transactionController.EXPECT().BeginTx(ctx).Return(tx, nil)

	balanceRepo.EXPECT().ChangeUserBalance(tx, -200, "test").Return(usecase.ErrInsufficientBalance)

	tx.EXPECT().Rollback().Return(port.ErrInsufficientBalance)

	err := userUsecase.Send(ctx, sendRequest)
	if !errors.Is(err, usecase.ErrInsufficientBalance) {
		t.Error(err)
	}
}

func TestSendWithNegativeAmount(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	shopRepo := repo.NewMockShopRepo(ctrl)
	balanceRepo := repo.NewMockUserBalanceRepo(ctrl)
	inventoryRepo := repo.NewMockUserInventoryRepo(ctrl)
	transactionController := repo.NewMockTransactionController(ctrl)
	userTransactionRepo := repo.NewMockUserTransactionRepo(ctrl)

	userUsecase := usecase.NewUser(shopRepo, balanceRepo, inventoryRepo, transactionController, userTransactionRepo)

	sendRequest := entity.SendCoinRequest{
		Amount:   -200,
		FromUser: "test",
		ToUser:   "test2",
	}

	ctx := context.Background()

	err := userUsecase.Send(ctx, sendRequest)
	if !errors.Is(err, usecase.ErrWrongCoinAmount) {
		t.Error(err)
	}
}
