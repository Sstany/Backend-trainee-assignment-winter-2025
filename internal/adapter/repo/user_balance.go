package repo

import (
	"context"
	"database/sql"
	"errors"

	"shop/internal/app/port"

	"go.uber.org/zap"
)

var _ port.UserBalanceRepo = (*Balance)(nil)

const (
	readBalance   = "SELECT balance FROM user_balance WHERE username=$1"
	changeBalance = "UPDATE balance FROM user_balance SET balance = balance + $1 WHERE username=$2"
)

type Balance struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewBalance(db *sql.DB, logger *zap.Logger) *Balance {
	return &Balance{
		db:     db,
		logger: logger,
	}
}

func (r *Balance) GetUserBalance(ctx context.Context, name string) (int, error) {
	var balance int

	err := r.db.QueryRowContext(ctx, readBalance, name).Scan(&balance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, port.ErrNotFound
		}
		return 0, err
	}

	return balance, nil
}

func (r *Balance) ChangeUserBalance(tx port.Transaction, count int, name string) error {
	tx.Exec(changeBalance)

	return nil
}
