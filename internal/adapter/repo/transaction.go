package repo

import (
	"context"
	"database/sql"

	"go.uber.org/zap"

	"shop/internal/app/port"
)

var _ port.TransactionController = (*TransactionSQL)(nil)

type TransactionSQL struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewTransactionSQL(db *sql.DB, logger *zap.Logger) *TransactionSQL {
	return &TransactionSQL{
		db:     db,
		logger: logger,
	}
}
func (r *TransactionSQL) BeginTx(ctx context.Context) (port.Transaction, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return tx, nil
}
