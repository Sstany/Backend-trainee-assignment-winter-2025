package usecase

import (
	"context"
	"errors"
	"fmt"

	"shop/internal/app/entity"
	"shop/internal/app/port"
)

var _ UserUseCase = (*User)(nil)

const (
	initCoins      = 1000
	minPasswordLen = 6
)

type User struct {
	authRepo              port.AuthRepo
	shopRepo              port.ShopRepo
	balanceRepo           port.UserBalanceRepo
	inventoryRepo         port.UserInventoryRepo
	passHasher            port.PassHasher
	transactionController port.TransactionController
	userTransactionRepo   port.UserTransactionRepo
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
	userTransactionRepo port.UserTransactionRepo,
) *User {
	return &User{
		authRepo:              authRepo,
		shopRepo:              shopRepo,
		balanceRepo:           balanceRepo,
		inventoryRepo:         inventoryRepo,
		passHasher:            passHasher,
		transactionController: transactionController,
		userTransactionRepo:   userTransactionRepo,
	}
}

func (r *User) Auth(ctx context.Context, login entity.Login) (*entity.Token, error) {
	if len(login.Password) < minPasswordLen {
		return nil, ErrUnsafePassword
	}

	if len(login.Password) > 20 {
		return nil, ErrLongPassword
	}

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

			if err = r.balanceRepo.Create(ctx, login.Username, initCoins); err != nil {
				return nil, fmt.Errorf("create user balance; %w", err)
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
		if errors.Is(err, port.ErrInsufficientBalance) {
			return ErrInsufficientBalance
		}
		return fmt.Errorf("change user balance: %w", err)
	}

	err = r.inventoryRepo.AddItem(tx, itemRequest.Username, itemRequest.Item)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("add item to user: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("buy trancastion: %w", err)
	}

	return nil
}

func (r *User) Send(ctx context.Context, sendReq entity.SendCoinRequest) error {
	if sendReq.Amount <= 0 {
		return ErrWrongCoinAmount
	}

	tx, err := r.transactionController.BeginTx(ctx)
	if err != nil {
		return err
	}

	err = r.balanceRepo.ChangeUserBalance(tx, -(sendReq.Amount), sendReq.FromUser)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, port.ErrInsufficientBalance) {
			return ErrInsufficientBalance
		}
		return fmt.Errorf("change user balance: %w", err)
	}

	err = r.balanceRepo.ChangeUserBalance(tx, sendReq.Amount, sendReq.ToUser)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("send coins to user: %w", err)
	}

	err = r.userTransactionRepo.SetUserTransaction(tx, sendReq)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("add to wallet transaction: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("buy trancastion: %w", err)
	}

	return nil
}
