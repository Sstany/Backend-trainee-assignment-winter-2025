package repo

import (
	"shop/internal/app/port"

	"go.uber.org/zap"
)

var _ port.ShopRepo = (*Shop)(nil)

type Shop struct {
	logger *zap.Logger
}

func NewShop(logger *zap.Logger) *Shop {
	return &Shop{
		logger: logger,
	}
}

func (r *Shop) GetItemPrice(name string) (int, bool) {
	var price int

	switch name {
	case "t-shirt":
		price = 80
	case "cup":
		price = 20
	case "book":
		price = 50
	case "pen":
		price = 10
	case "powerbank":
		price = 200
	case "hoody":
		price = 300
	case "umbrella":
		price = 200
	case "socks":
		price = 10
	case "wallet":
		price = 50
	case "pink-hoody":
		price = 500

	default:
		return 0, false
	}

	return price, true
}
