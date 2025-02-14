package repo

import (
	"context"
	"database/sql"
	"fmt"

	"go.uber.org/zap"

	"shop/internal/app/entity"
	"shop/internal/app/port"
)

var _ port.UserInventoryRepo = (*Inventory)(nil)

const (
	addItem  = "INSERT INTO inventory (username, item) VALUES($1, $2)"
	getItems = "SELECT item, COUNT(*) AS count FROM inventory item WHERE username=$1 GROUP BY item"
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

func (r *Inventory) Get(ctx context.Context, username string) ([]entity.Inventory, error) {
	rows, err := r.db.QueryContext(ctx, getItems, username)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var result []entity.Inventory
	var temp entity.Inventory

	for rows.Next() {
		if err = rows.Scan(&temp.Type, &temp.Quantity); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		result = append(result, temp)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows err: %w", err)
	}

	return result, nil
}
