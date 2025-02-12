package port

import (
	"context"
)

//go:generate mockgen -destination ../../adapter/repo/user_balance_mock.go -package repo -source ./user_balance.go

type UserBalanceRepo interface {
	GetUserBalance(ctx context.Context, name string) (int, error)
	ChangeUserBalance(tx Transaction, count int, name string) error
}
