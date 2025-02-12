package port

import (
	"context"
	"database/sql"
)

var _ Transaction = (*sql.Tx)(nil)

//go:generate mockgen -destination ../../adapter/repo/mock/transaction_mock.go -package repo -source ./transaction.go

type TransactionController interface {
	BeginTx(ctx context.Context) (Transaction, error)
}

type Transaction interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	Commit() error
	Rollback() error
}
