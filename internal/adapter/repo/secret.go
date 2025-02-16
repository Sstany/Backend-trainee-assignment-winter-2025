package repo

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"time"

	"shop/internal/app/entity"
	"shop/internal/app/port"
	"shop/internal/pkg"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var _ port.SecretRepo = (*Secret)(nil)
var ErrInvalidSigningAlgo = errors.New("invalid signing algo")

const defaultBits = 2048

const defaultAccessExpiration = time.Hour * 24

type Secret struct {
	privateKey      *rsa.PrivateKey
	publicKey       *rsa.PublicKey
	jwtIssuer       string
	signingAlgoName string
	logger          *zap.Logger
}

func NewSecret(logger *zap.Logger, jwtSigningKeyPath string, jwtIssuer string) *Secret {
	var key *rsa.PrivateKey

	info, err := os.Stat(jwtSigningKeyPath)
	if os.IsNotExist(err) || info.IsDir() {
		key, err = rsa.GenerateKey(rand.Reader, defaultBits)
		if err != nil {
			logger.Error("generate key failed", zap.Error(err))
			panic(err)
		}

		keyBytes := x509.MarshalPKCS1PrivateKey(key)
		pemBlock := &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: keyBytes,
		}

		//nolint:govet // ok.
		file, err := os.Create(jwtSigningKeyPath)
		if err != nil {
			logger.Error("create key file failed", zap.Error(err))
			panic(err)
		}

		defer file.Close()
		err = pem.Encode(file, pemBlock)
		if err != nil {
			logger.Error("encode key failed", zap.Error(err))
			panic(err)
		}
	} else {
		//nolint:govet // ok.
		privateKey, err := os.ReadFile(jwtSigningKeyPath)
		if err != nil {
			logger.Error("read key failed", zap.Error(err))
			panic(err)
		}

		block, _ := pem.Decode(privateKey)

		key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			logger.Error("private key parse failed", zap.Error(err))
			panic(err)
		}
	}

	return &Secret{
		privateKey:      key,
		publicKey:       &key.PublicKey,
		jwtIssuer:       jwtIssuer,
		signingAlgoName: jwt.SigningMethodRS256.Name,
		logger:          logger,
	}
}

func (r *Secret) CreateToken(username string) (*entity.Token, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256,
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

	return pkg.PointerTo(entity.Token(accessToken)), nil
}

func (r *Secret) ParseJWT(token string) (jwt.MapClaims, error) {
	var mapClaims jwt.MapClaims

	_, err := jwt.ParseWithClaims(token, &mapClaims, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != r.signingAlgoName {
			return nil, ErrInvalidSigningAlgo
		}

		return r.publicKey, nil
	})
	if err != nil {
		return nil, err
	}
	return mapClaims, nil
}
