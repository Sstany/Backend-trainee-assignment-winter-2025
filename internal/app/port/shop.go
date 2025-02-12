package port

import "context"

//go:generate mockgen -destination ../../adapter/repo/shop_mock.go -package repo -source ./shop.go

type ShopRepo interface {
	GetItemPrice(ctx context.Context, name string) (int, bool)
}
