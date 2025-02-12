package usecase

import (
	"context"
	"errors"
	"fmt"

	"shop/internal/app/entity"
	"shop/internal/app/port"
)

var _ UserUseCase = (*User)(nil)

type User struct {
	authRepo   port.AuthRepo
	passHasher port.PassHasher
	// userInfoRepo
	// transactionRepo
	// inventoryRepo
}

func NewUser(authRepo port.AuthRepo, passHasher port.PassHasher) *User {
	return &User{
		authRepo:   authRepo,
		passHasher: passHasher,
	}
}

func (r *User) Auth(ctx context.Context, login entity.Login) (*entity.Token, error) {
	passHash, err := r.authRepo.ReadPassword(ctx, login.Username)
	if err != nil {
		if errors.Is(err, port.ErrNotFound) {
			hash, err := r.passHasher.Hash(login.Password)
			if err != nil {
				return nil, fmt.Errorf("hash password: %w", err)
			}

			login.Password = hash

			if err = r.authRepo.CreateUser(ctx, login); err != nil {
				return nil, fmt.Errorf("create user: %w", err)
			}

			// TODO: return jwt

			return nil, nil
		}

		return nil, fmt.Errorf("read password: %w", err)
	}

	if passHash != "" && !r.passHasher.Compare(login.Password, passHash) {
		return nil, ErrInvalidPassword
	}

	// TODO: create jwt with username here.

	return nil, nil
}
