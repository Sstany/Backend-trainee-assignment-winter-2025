package repo

import (
	"go.uber.org/zap"

	"shop/internal/app/port"
)

var _ port.UserInventoryRepo = (*Inventory)(nil)

const (
	addItem = "INSERT INTO inventory username, itemname VALUES($1, $2)"
)

type Inventory struct {
	logger *zap.Logger
}

func NewInventory(logger *zap.Logger) *Inventory {
	return &Inventory{
		logger: logger,
	}
}
func (r *Inventory) AddItem(tx port.Transaction, username string, item string) error {
	_, err := tx.Exec(addItem)
	if err != nil {
		return err
	}

	return nil
}
