package usecase_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"

	"shop/internal/adapter/password"
	repo "shop/internal/adapter/repo"
	"shop/internal/app/entity"
	"shop/internal/app/usecase"
)

func TestBuy(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	authRepo := repo.NewMockAuthRepo(ctrl)

	passHasher := password.NewMockPassHasher(ctrl)

	shopRepo := repo.NewMockShopRepo(ctrl)
	balanceRepo := repo.NewMockUserBalanceRepo(ctrl)
	inventoryRepo := repo.NewMockUserInventoryRepo(ctrl)
	transactionController := repo.NewMockTransactionController(ctrl)

	userUsecase := usecase.NewUser(authRepo, shopRepo, balanceRepo, inventoryRepo, passHasher, transactionController)

	itemRequest := entity.ItemRequest{
		Username: "test",
		Item:     "hoody",
	}
	ctx := context.Background()
	tx := repo.NewMockTransaction(ctrl)

	shopRepo.EXPECT().GetItemPrice(ctx, "hoody").Return(300, true)
	transactionController.EXPECT().BeginTx(ctx).Return(tx, nil)
	balanceRepo.EXPECT().ChangeUserBalance(tx, -300, "test").Return(nil)
	inventoryRepo.EXPECT().AddItem(tx, "test", "hoody").Return(nil)
	tx.EXPECT().Commit().Return(nil)

	err := userUsecase.Buy(ctx, itemRequest)
	if err != nil {
		t.Error(err)
	}
	// TODO: Check that JWT is created with right username
}
