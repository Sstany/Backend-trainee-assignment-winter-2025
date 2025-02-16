package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"shop/internal/app/entity"
	"shop/internal/app/port"

	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
)

const (
	maxPasswordLen = 20
)

type Auth struct {
	authRepo    port.AuthRepo
	balanceRepo port.UserBalanceRepo
	passHasher  port.PassHasher
	secretRepo  port.SecretRepo

	logger *zap.Logger
}

func NewAuth(
	authRepo port.AuthRepo,
	balanceRepo port.UserBalanceRepo,
	passHasher port.PassHasher,
	secretRepo port.SecretRepo,
	logger *zap.Logger,
) (*Auth, error) {
	as := &Auth{
		authRepo:    authRepo,
		balanceRepo: balanceRepo,
		passHasher:  passHasher,
		secretRepo:  secretRepo,
		logger:      logger,
	}

	return as, nil
}

func (r *Auth) Auth(ctx context.Context, login entity.Login) (*entity.Token, error) {
	if len(login.Password) < minPasswordLen {
		return nil, ErrUnsafePassword
	}

	if len(login.Password) > maxPasswordLen {
		return nil, ErrLongPassword
	}

	passHash, err := r.authRepo.ReadPassword(ctx, login.Username)
	if err != nil {
		// User registration.
		//nolint:govet // ok.
		if errors.Is(err, port.ErrNotFound) {
			hash, err := r.passHasher.Hash(login.Password)
			if err != nil {
				return nil, fmt.Errorf("hash password: %w", err)
			}

			login.Password = hash

			if err = r.authRepo.CreateUser(ctx, login); err != nil {
				return nil, fmt.Errorf("create user: %w", err)
			}

			if err = r.balanceRepo.Create(ctx, login.Username, initCoins); err != nil {
				return nil, fmt.Errorf("create user balance; %w", err)
			}

			return r.secretRepo.CreateToken(login.Username)
		}

		return nil, fmt.Errorf("read password: %w", err)
	}

	if passHash != "" && !r.passHasher.Compare(login.Password, passHash) {
		return nil, ErrInvalidPassword
	}

	return r.secretRepo.CreateToken(login.Username)
}

func (r *Auth) AuthenticateWithAccessToken(token string) (string, bool, error) {
	mapClaims, err := r.secretRepo.ParseJWT(token)

	var username string

	usernameRaw, ok := mapClaims["uid"]
	if ok {
		username, ok = usernameRaw.(string)
		if !ok {
			return "", false, ErrInvalidToken
		}
	}

	if err != nil {
		var vErr *jwt.ValidationError

		if errors.As(err, &vErr) {
			if vErr.Errors == jwt.ValidationErrorExpired {
				return username, false, ErrTokenExpired
			}
		}

		return username, false, err
	}

	if !mapClaims.VerifyExpiresAt(time.Now().Unix(), true) {
		return username, false, ErrTokenExpired
	}

	return username, true, nil
}
