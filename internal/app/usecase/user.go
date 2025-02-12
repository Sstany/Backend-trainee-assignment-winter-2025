package usecase

import (
	"context"
	"errors"
	"fmt"

	"shop/internal/app/entity"
	"shop/internal/app/port"
)

var _ UserUseCase = (*User)(nil)

type User struct {
	authRepo              port.AuthRepo
	shopRepo              port.ShopRepo
	balanceRepo           port.UserBalanceRepo
	inventoryRepo         port.UserInventoryRepo
	passHasher            port.PassHasher
	transactionController port.TransactionController
	// userInfoRepo
	// transactionRepos
	// inventoryRepo
}

func NewUser(
	authRepo port.AuthRepo,
	shopRepo port.ShopRepo,
	balanceRepo port.UserBalanceRepo,
	inventoryRepo port.UserInventoryRepo,
	passHasher port.PassHasher,
	transactionController port.TransactionController,
) *User {
	return &User{
		authRepo:              authRepo,
		shopRepo:              shopRepo,
		balanceRepo:           balanceRepo,
		inventoryRepo:         inventoryRepo,
		passHasher:            passHasher,
		transactionController: transactionController,
	}
}

func (r *User) Auth(ctx context.Context, login entity.Login) (*entity.Token, error) {
	passHash, err := r.authRepo.ReadPassword(ctx, login.Username)
	if err != nil {
		if errors.Is(err, port.ErrNotFound) {
			hash, err := r.passHasher.Hash(login.Password)
			if err != nil {
				return nil, fmt.Errorf("hash password: %w", err)
			}

			login.Password = hash

			if err = r.authRepo.CreateUser(ctx, login); err != nil {
				return nil, fmt.Errorf("create user: %w", err)
			}

			// TODO: return jwt

			return nil, nil
		}

		return nil, fmt.Errorf("read password: %w", err)
	}

	if passHash != "" && !r.passHasher.Compare(login.Password, passHash) {
		return nil, ErrInvalidPassword
	}

	// TODO: create jwt with username here.

	return nil, nil
}

func (r *User) Buy(ctx context.Context, itemRequest entity.ItemRequest) error {
	price, exists := r.shopRepo.GetItemPrice(ctx, itemRequest.Item)
	if !exists {
		return ErrItemNotExists
	}

	tx, err := r.transactionController.BeginTx(ctx)
	if err != nil {
		return err
	}

	err = r.balanceRepo.ChangeUserBalance(tx, -price, itemRequest.Username)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = r.inventoryRepo.AddItem(tx, itemRequest.Username, itemRequest.Item)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
