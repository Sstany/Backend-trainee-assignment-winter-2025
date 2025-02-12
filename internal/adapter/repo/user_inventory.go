package repo

import (
	"database/sql"

	"go.uber.org/zap"

	"shop/internal/app/port"
)

var _ port.UserInventoryRepo = (*Inventory)(nil)

const (
	addItem = "INSERT INTO inventory (username, item) VALUES($1, $2)"
)

type Inventory struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewInventory(db *sql.DB, logger *zap.Logger) *Inventory {
	return &Inventory{
		db:     db,
		logger: logger,
	}
}

func (r *Inventory) AddItem(tx port.Transaction, username string, item string) error {
	_, err := tx.Exec(addItem, username, item)
	if err != nil {
		return err
	}

	return nil
}
