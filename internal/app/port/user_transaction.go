package port

import "shop/internal/app/entity"

//go:generate mockgen -destination ../../adapter/repo/mock/user_transaction_mock.go -package repo -source ./user_transaction.go

type UserTransactionRepo interface {
	SetUserTransaction(tx Transaction, sendCoin entity.SendCoinRequest) error
}
