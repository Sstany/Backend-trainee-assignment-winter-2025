package port

//go:generate mockgen -destination ../../adapter/repo/mock/shop_mock.go -package repo -source ./shop.go

type ShopRepo interface {
	GetItemPrice(name string) (int, bool)
}
