package usecase

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"

	"shop/internal/app/entity"
	"shop/internal/app/port"
)

var _ UserUseCase = (*User)(nil)

const (
	initCoins      = 1000
	minPasswordLen = 6
)

type User struct {
	shopRepo              port.ShopRepo
	balanceRepo           port.UserBalanceRepo
	inventoryRepo         port.UserInventoryRepo
	transactionController port.TransactionController
	userTransactionRepo   port.UserTransactionRepo
	logger                *zap.Logger
}

func NewUser(
	shopRepo port.ShopRepo,
	balanceRepo port.UserBalanceRepo,
	inventoryRepo port.UserInventoryRepo,
	transactionController port.TransactionController,
	userTransactionRepo port.UserTransactionRepo,
	logger *zap.Logger,
) *User {
	return &User{
		shopRepo:              shopRepo,
		balanceRepo:           balanceRepo,
		inventoryRepo:         inventoryRepo,
		transactionController: transactionController,
		userTransactionRepo:   userTransactionRepo,
		logger:                logger,
	}
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
		err := tx.Rollback()
		if err != nil {
			r.logger.Error("transaction failed", zap.Error(err))
		}

		if errors.Is(err, port.ErrInsufficientBalance) {
			return ErrInsufficientBalance
		}
		if errors.Is(err, port.ErrReveicerNotExists) {
			return ErrReveicerNotExists
		}

		return fmt.Errorf("change user balance: %w", err)
	}

	err = r.inventoryRepo.AddItem(tx, itemRequest.Username, itemRequest.Item)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			r.logger.Error("transaction failed", zap.Error(err))
		}
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
		err := tx.Rollback()
		if err != nil {
			r.logger.Error("transaction failed", zap.Error(err))
		}

		return fmt.Errorf("change user balance: %w", err)
	}

	err = r.balanceRepo.ChangeUserBalance(tx, sendReq.Amount, sendReq.ToUser)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			r.logger.Error("transaction failed", zap.Error(err))
		}

		return fmt.Errorf("send coins to user: %w", err)
	}

	err = r.userTransactionRepo.SetUserTransaction(tx, sendReq)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			r.logger.Error("transaction failed", zap.Error(err))
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("buy trancastion: %w", err)
	}

	return nil
}

func (r *User) Info(ctx context.Context, infoReq entity.InfoRequest) (*entity.InfoResponse, error) {
	items, err := r.inventoryRepo.Get(ctx, infoReq.Username)
	if err != nil {
		return nil, fmt.Errorf("get items with count; %w", err)
	}

	coins, err := r.balanceRepo.GetUserBalance(ctx, infoReq.Username)
	if err != nil {
		return nil, fmt.Errorf("user balance: %w", err)
	}

	rec, err := r.userTransactionRepo.GetRecievedOperations(ctx, infoReq.Username)
	if err != nil {
		return nil, fmt.Errorf("user recieved operations: %w", err)
	}

	sent, err := r.userTransactionRepo.GetSentOperations(ctx, infoReq.Username)
	if err != nil {
		return nil, fmt.Errorf("user sent operations: %w", err)
	}

	info := &entity.InfoResponse{
		CoinHistory: entity.CoinHistory{
			Received: rec,
			Sent:     sent,
		},
		Coins:     coins,
		Inventory: items,
	}

	return info, err
}
