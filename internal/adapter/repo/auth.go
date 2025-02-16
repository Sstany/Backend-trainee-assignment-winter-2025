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

var _ port.AuthRepo = (*Auth)(nil)

const (
	readPasswordByUsername = `SELECT password FROM users WHERE users.username = $1`
	createUser             = "INSERT INTO users (username, password) VALUES ($1, $2)"
)

const codeUniqueViolation = "23505"

type Auth struct {
	db                         *sql.DB
	stmtReadPasswordByUsername *sql.Stmt
	stmtCreateUser             *sql.Stmt
	logger                     *zap.Logger
}

func NewAuth(db *sql.DB, logger *zap.Logger) (*Auth, error) {
	readPassStmt, err := db.Prepare(readPasswordByUsername)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}

	createUserStmt, err := db.Prepare(createUser)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}

	return &Auth{
		db:                         db,
		stmtReadPasswordByUsername: readPassStmt,
		stmtCreateUser:             createUserStmt,
		logger:                     logger,
	}, nil
}

func (r *Auth) ReadPassword(ctx context.Context, username string) (string, error) {
	var passHash string

	err := r.stmtReadPasswordByUsername.QueryRowContext(
		ctx,
		username,
	).Scan(&passHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", port.ErrNotFound
		}

		return "", err
	}

	return passHash, nil
}

func (r *Auth) CreateUser(ctx context.Context, login entity.Login) error {
	_, err := r.stmtCreateUser.ExecContext(ctx, login.Username, login.Password)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == codeUniqueViolation {
				return port.ErrAlreadyRegistred
			}
		}
		return err
	}

	return nil
}
