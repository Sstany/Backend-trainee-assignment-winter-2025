package repo

import (
	"database/sql"

	"shop/internal/app/entity"
	"shop/internal/app/port"

	"go.uber.org/zap"
)

var _ port.UserTransactionRepo = (*UserTransaction)(nil)

const (
	setUserTransaction = "INSERT INTO user_transactions (from, to, count) VALUES($1, $2, $3)"
)

type UserTransaction struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewUserTransaction(db *sql.DB, logger *zap.Logger) *UserTransaction {
	return &UserTransaction{
		db:     db,
		logger: logger,
	}
}

func (*UserTransaction) SetUserTransaction(tx port.Transaction, sendCoin entity.SendCoinRequest) error {
	_, err := tx.Exec(setUserTransaction, sendCoin.FromUser, sendCoin.ToUser, sendCoin.Amount)
	if err != nil {
		return err
	}

	return nil
}
