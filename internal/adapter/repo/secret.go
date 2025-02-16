package repo

import (
	"crypto/rand"
	"errors"
	"fmt"
	"os"
	"time"

	"shop/internal/app/entity"
	"shop/internal/app/port"
	"shop/internal/pkg"

	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
)

var _ port.SecretRepo = (*Secret)(nil)

var ErrInvalidSigningAlgo = errors.New("invalid signing algo")

const defaultBits = 2048

const defaultAccessExpiration = time.Hour * 24

type Secret struct {
	secretKey       []byte
	signingAlgoName string
	logger          *zap.Logger
}

func NewSecret(logger *zap.Logger, jwtSigningKeyPath string, jwtIssuer string) *Secret {
	var key []byte

	info, err := os.Stat(jwtSigningKeyPath)
	if os.IsNotExist(err) || info.IsDir() {
		key, err = generateHMACKey()
		if err != nil {
			logger.Error("generate key failed", zap.Error(err))
			panic(err)
		}

		//nolint:govet // ok.
		file, err := os.Create(jwtSigningKeyPath)
		if err != nil {
			logger.Error("create key file failed", zap.Error(err))
			panic(err)
		}

		defer file.Close()
	} else {
		//nolint:govet // ok.
		fileKey, err := os.ReadFile(jwtSigningKeyPath)
		if err != nil {
			logger.Error("read key failed", zap.Error(err))
			panic(err)
		}

		key = fileKey
	}

	return &Secret{
		secretKey:       key,
		signingAlgoName: jwt.SigningMethodHS256.Name,
		logger:          logger,
	}
}

func (r *Secret) CreateToken(username string) (*entity.Token, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"exp": time.Now().Add(defaultAccessExpiration).Unix(),
			"uid": username,
		},
	)

	accessToken, err := token.SignedString(r.secretKey)
	if err != nil {
		return nil, fmt.Errorf("jwt: %w", err)
	}

	return pkg.PointerTo(entity.Token(accessToken)), nil
}

func (r *Secret) ParseJWT(token string) (jwt.MapClaims, error) {
	var mapClaims jwt.MapClaims

	_, err := jwt.ParseWithClaims(token, &mapClaims, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != r.signingAlgoName {
			return nil, ErrInvalidSigningAlgo
		}

		return r.secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	return mapClaims, nil
}

func generateHMACKey() ([]byte, error) {
	key := make([]byte, 32) // 256-bit key
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}
