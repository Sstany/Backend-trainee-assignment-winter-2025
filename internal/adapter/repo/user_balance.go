package repo

import (
	"context"
	"database/sql"
	"errors"

	"shop/internal/app/port"

	"github.com/lib/pq"
	"go.uber.org/zap"
)

var _ port.UserBalanceRepo = (*Balance)(nil)

const (
	readBalance   = "SELECT balance FROM wallets WHERE username=$1"
	changeBalance = "UPDATE wallets SET balance = balance + $1 WHERE username=$2"
	create        = "INSERT INTO wallets (username, balance) VALUES($1, $2)"
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
	res, err := tx.Exec(changeBalance, count, name)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23514" {
				return port.ErrInsufficientBalance
			}
		}
		return err
	}
	if str, _ := res.RowsAffected(); str != 1 {
		return port.ErrReveicerNotExists
	}

	return nil
}

func (r *Balance) Create(ctx context.Context, name string, initCoins int) error {
	_, err := r.db.ExecContext(ctx, create, name, initCoins)
	return err
}
