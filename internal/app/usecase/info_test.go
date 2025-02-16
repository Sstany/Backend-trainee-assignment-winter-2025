package usecase_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"

	repo "shop/internal/adapter/repo/mock"
	"shop/internal/app/entity"
	"shop/internal/app/usecase"
)

func TestInfo(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	shopRepo := repo.NewMockShopRepo(ctrl)
	balanceRepo := repo.NewMockUserBalanceRepo(ctrl)
	inventoryRepo := repo.NewMockUserInventoryRepo(ctrl)
	transactionController := repo.NewMockTransactionController(ctrl)
	userTransactionRepo := repo.NewMockUserTransactionRepo(ctrl)

	userUsecase := usecase.NewUser(shopRepo, balanceRepo, inventoryRepo, transactionController, userTransactionRepo, nil)

	request := entity.InfoRequest{
		Username: "test",
	}

	ctx := context.Background()

	inventoryRepo.EXPECT().Get(ctx, "test").Return([]entity.Inventory{}, nil)
	balanceRepo.EXPECT().GetUserBalance(ctx, "test").Return(3, nil)
	userTransactionRepo.EXPECT().GetRecievedOperations(ctx, "test").Return([]entity.Received{}, nil)
	userTransactionRepo.EXPECT().GetSentOperations(ctx, "test").Return([]entity.Sent{}, nil)

	_, err := userUsecase.Info(ctx, request)
	if err != nil {
		t.Error(err)
	}
}
