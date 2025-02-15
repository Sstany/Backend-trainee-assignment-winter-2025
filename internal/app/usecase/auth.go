package usecase

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"time"

	"shop/internal/app/entity"
	"shop/internal/app/port"
	"shop/internal/pkg"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwk"
	"go.uber.org/zap"
)

const (
	defaultAccessExpiration = time.Hour * 24
	maxPasswordLen          = 20
)

type Auth struct {
	authRepo    port.AuthRepo
	balanceRepo port.UserBalanceRepo
	passHasher  port.PassHasher

	logger *zap.Logger

	privateKey      *ecdsa.PrivateKey
	publicKey       *ecdsa.PublicKey
	jwkPublicKey    jwk.ECDSAPublicKey
	jwtIssuer       string
	signingAlgoName string
}

func NewAuth(
	authRepo port.AuthRepo,
	balanceRepo port.UserBalanceRepo,
	passHasher port.PassHasher,
	secretRepo port.SecretRepo,
	logger *zap.Logger,
) (*Auth, error) {
	as := &Auth{
		authRepo:        authRepo,
		balanceRepo:     balanceRepo,
		passHasher:      passHasher,
		signingAlgoName: jwt.SigningMethodES512.Name,
		logger:          logger.Named("authService"),
	}

	as.privateKey = secretRepo.PrivateKey()
	as.publicKey = secretRepo.PublicKey()
	as.jwtIssuer = secretRepo.JWTIssuer()

	as.jwkPublicKey = jwk.NewECDSAPublicKey()
	if err := as.jwkPublicKey.FromRaw(as.publicKey); err != nil {
		panic(err)
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

			return r.createToken(login.Username)
		}

		return nil, fmt.Errorf("read password: %w", err)
	}

	if passHash != "" && !r.passHasher.Compare(login.Password, passHash) {
		return nil, ErrInvalidPassword
	}

	return r.createToken(login.Username)
}

func (r *Auth) createToken(username string) (*entity.Token, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES512,
		jwt.MapClaims{
			"jti": uuid.NewString(),
			"iss": r.jwtIssuer,
			"exp": time.Now().Add(defaultAccessExpiration).Unix(),
			"uid": username,
		},
	)

	accessToken, err := token.SignedString(r.privateKey)
	if err != nil {
		return nil, fmt.Errorf("jwt:s %w", err)
	}

	return pkg.PointerTo[entity.Token](entity.Token(accessToken)), nil
}

func (r *Auth) AuthenticateWithAccessToken(token string) (string, bool, error) {
	var mapClaims jwt.MapClaims

	_, err := jwt.ParseWithClaims(token, &mapClaims, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != r.signingAlgoName {
			return nil, ErrInvalidSigningAlgo
		}

		return r.publicKey, nil
	})

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
