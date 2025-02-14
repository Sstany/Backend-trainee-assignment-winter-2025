package port

import (
	"context"

	"shop/internal/app/entity"
)

//go:generate mockgen -destination ../../adapter/repo/mock/user_inventory_mock.go -package repo -source ./user_inventory.go

type UserInventoryRepo interface {
	AddItem(tx Transaction, username string, item string) error
	Get(ctx context.Context, username string) ([]entity.Inventory, error)
}
