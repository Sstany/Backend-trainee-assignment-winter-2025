package repo

import (
	"context"
	"database/sql"
	"errors"

	"go.uber.org/zap"

	"shop/internal/app/entity"
	"shop/internal/app/port"
)

var _ port.AuthRepo = (*Auth)(nil)

const (
	readPasswordByUsername = `SELECT password FROM users WHERE users.username = $1`
	createUser             = "INSERT INTO users (username, password) VALUES ($1, $2)"
)

type Auth struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewAuth(db *sql.DB, logger *zap.Logger) *Auth {
	return &Auth{
		db:     db,
		logger: logger,
	}
}

func (r *Auth) ReadPassword(ctx context.Context, username string) (string, error) {
	var passHash string

	if err := r.db.QueryRowContext(ctx, readPasswordByUsername, username).Scan(&passHash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", port.ErrNotFound
		}

		return "", err
	}

	return passHash, nil
}

func (r *Auth) CreateUser(ctx context.Context, login entity.Login) error {
	_, err := r.db.ExecContext(ctx, createUser, login.Username, login.Password)
	return err
}
