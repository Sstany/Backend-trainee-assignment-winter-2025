package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
	"go.uber.org/zap"

	"shop/internal/app/entity"
	"shop/internal/app/port"
)

var _ port.UserInventoryRepo = (*Inventory)(nil)

const (
	addItem  = "INSERT INTO inventory (username, item) VALUES($1, $2)"
	getItems = "SELECT item, COUNT(*) AS count FROM inventory item WHERE username=$1 GROUP BY item"
)

const codeSerializationFailure = "40001"

type Inventory struct {
	db           *sql.DB
	stmtGetItems *sql.Stmt
	logger       *zap.Logger
}

func NewInventory(db *sql.DB, logger *zap.Logger) (*Inventory, error) {
	getItemsStmt, err := db.Prepare(getItems)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare getItems statement: %w", err)
	}

	return &Inventory{
		db:           db,
		stmtGetItems: getItemsStmt,
		logger:       logger,
	}, nil
}

func (r *Inventory) AddItem(tx port.Transaction, username string, item string) error {
	_, err := tx.Exec(addItem, username, item)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == codeSerializationFailure {
				return port.ErrTransactionFailure
			}
		}

		return err
	}

	return nil
}

func (r *Inventory) Get(ctx context.Context, username string) ([]entity.Inventory, error) {
	rows, err := r.stmtGetItems.QueryContext(ctx, username)
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
