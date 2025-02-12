package port

//go:generate mockgen -destination ../../adapter/repo/mock/user_inventory_mock.go -package repo -source ./user_inventory.go

type UserInventoryRepo interface {
	AddItem(tx Transaction, userName string, item string) error
}
