package port

import (
	"context"

	"shop/internal/app/entity"
)

//go:generate mockgen -destination ../../adapter/repo/mock/user_transaction_mock.go -package repo -source ./user_transaction.go

type UserTransactionRepo interface {
	SetUserTransaction(tx Transaction, sendCoin entity.SendCoinRequest) error
	GetRecievedOperations(ctx context.Context, username string) ([]entity.Received, error)
	GetSentOperations(ctx context.Context, username string) ([]entity.Sent, error)
}
