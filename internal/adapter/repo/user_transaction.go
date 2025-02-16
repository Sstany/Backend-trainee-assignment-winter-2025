package repo

import (
	"context"
	"database/sql"
	"fmt"

	"shop/internal/app/entity"
	"shop/internal/app/port"

	"go.uber.org/zap"
)

var _ port.UserTransactionRepo = (*UserTransaction)(nil)

const (
	setUserTransaction = "INSERT INTO user_transactions (user_from, user_to, amount) VALUES($1, $2, $3)"
	getRecieved        = "SELECT user_from, amount FROM user_transactions WHERE user_to = $1"
	getSent            = "SELECT user_to, amount FROM user_transactions WHERE user_from = $1"
)

type UserTransaction struct {
	db              *sql.DB
	stmtGetRecieved *sql.Stmt
	stmtGetSent     *sql.Stmt
	logger          *zap.Logger
}

func NewUserTransaction(db *sql.DB, logger *zap.Logger) (*UserTransaction, error) {
	getRecievedStmt, err := db.Prepare(getRecieved)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare getRecieved statement: %w", err)
	}

	getSentStmt, err := db.Prepare(getSent)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare getSent statement: %w", err)
	}

	return &UserTransaction{
		db:              db,
		stmtGetRecieved: getRecievedStmt,
		stmtGetSent:     getSentStmt,
		logger:          logger,
	}, nil
}

func (r *UserTransaction) SetUserTransaction(tx port.Transaction, sendCoin entity.SendCoinRequest) error {
	_, err := tx.Exec(setUserTransaction, sendCoin.FromUser, sendCoin.ToUser, sendCoin.Amount)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserTransaction) GetRecievedOperations(ctx context.Context, username string) ([]entity.Received, error) {
	rows, err := r.stmtGetRecieved.QueryContext(ctx, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recieved entity.Received
	var res []entity.Received

	for rows.Next() {
		if err = rows.Scan(
			&recieved.FromUser,
			&recieved.Amount,
		); err != nil {
			r.logger.Error("get user recieved operations failed", zap.String("username", username), zap.Error(err))
			break
		}
		res = append(res, recieved)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (r *UserTransaction) GetSentOperations(ctx context.Context, username string) ([]entity.Sent, error) {
	rows, err := r.stmtGetSent.QueryContext(ctx, username)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var sent entity.Sent
	var res []entity.Sent

	for rows.Next() {
		if err = rows.Scan(
			&sent.ToUser,
			&sent.Amount,
		); err != nil {
			r.logger.Error("get user sent operations failed", zap.String("username", username), zap.Error(err))
			break
		}
		res = append(res, sent)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}
