package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"go.uber.org/zap"

	repo "shop/internal/adapter/repo/mock"

	"shop/internal/app/entity"
	"shop/internal/app/port"
	"shop/internal/app/usecase"
)

func TestBuy(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	shopRepo := repo.NewMockShopRepo(ctrl)
	balanceRepo := repo.NewMockUserBalanceRepo(ctrl)

	inventoryRepo := repo.NewMockUserInventoryRepo(ctrl)
	transactionController := repo.NewMockTransactionController(ctrl)
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Error(err)
	}

	userUsecase := usecase.NewUser(shopRepo, balanceRepo, inventoryRepo, transactionController, nil, logger)

	itemRequest := entity.ItemRequest{
		Username: "test",
		Item:     "hoody",
	}
	ctx := context.Background()
	tx := repo.NewMockTransaction(ctrl)

	shopRepo.EXPECT().GetItemPrice("hoody").Return(300, true)
	transactionController.EXPECT().BeginTx(ctx).Return(tx, nil)
	balanceRepo.EXPECT().ChangeUserBalance(tx, -300, "test").Return(nil)
	inventoryRepo.EXPECT().AddItem(tx, "test", "hoody").Return(nil)
	tx.EXPECT().Commit().Return(nil)

	err = userUsecase.Buy(ctx, itemRequest)
	if err != nil {
		t.Error(err)
	}
}

func TestBuyNotExistItem(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	shopRepo := repo.NewMockShopRepo(ctrl)
	balanceRepo := repo.NewMockUserBalanceRepo(ctrl)
	inventoryRepo := repo.NewMockUserInventoryRepo(ctrl)
	transactionController := repo.NewMockTransactionController(ctrl)
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Error(err)
	}

	userUsecase := usecase.NewUser(shopRepo, balanceRepo, inventoryRepo, transactionController, nil, logger)

	itemRequest := entity.ItemRequest{
		Username: "test",
		Item:     "something",
	}
	ctx := context.Background()

	shopRepo.EXPECT().GetItemPrice("something").Return(0, false)

	err = userUsecase.Buy(ctx, itemRequest)
	if !errors.Is(err, usecase.ErrItemNotExists) {
		t.Error(err)
	}
}

func TestWithInsufficientBalance(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	shopRepo := repo.NewMockShopRepo(ctrl)
	balanceRepo := repo.NewMockUserBalanceRepo(ctrl)
	inventoryRepo := repo.NewMockUserInventoryRepo(ctrl)
	transactionController := repo.NewMockTransactionController(ctrl)
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Error(err)
	}

	userUsecase := usecase.NewUser(shopRepo, balanceRepo, inventoryRepo, transactionController, nil, logger)

	itemRequest := entity.ItemRequest{
		Username: "test",
		Item:     "hoody",
	}
	ctx := context.Background()
	tx := repo.NewMockTransaction(ctrl)

	shopRepo.EXPECT().GetItemPrice("hoody").Return(300, true)
	transactionController.EXPECT().BeginTx(ctx).Return(tx, nil)
	balanceRepo.EXPECT().ChangeUserBalance(tx, -300, "test").Return(port.ErrInsufficientBalance)
	tx.EXPECT().Rollback().Return(nil)

	err = userUsecase.Buy(ctx, itemRequest)
	if !errors.Is(err, usecase.ErrInsufficientBalance) {
		t.Error(err)
	}
}
