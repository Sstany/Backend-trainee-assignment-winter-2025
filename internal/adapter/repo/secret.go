package repo

import (
	"crypto/ecdsa"
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

const defaultAccessExpiration = time.Hour * 24

type Secret struct {
	privateKey      *ecdsa.PrivateKey
	publicKey       *ecdsa.PublicKey
	jwtIssuer       string
	signingAlgoName string
	logger          *zap.Logger
}

func (r *Secret) PrivateKey() *ecdsa.PrivateKey { return r.privateKey }
func (r *Secret) PublicKey() *ecdsa.PublicKey   { return r.publicKey }
func (r *Secret) JWTIssuer() string             { return r.jwtIssuer }

func NewSecret(logger *zap.Logger, jwtSigningKeyPath string, jwtIssuer string) *Secret {
	logger = logger.Named("secret")

	privateKey, err := os.ReadFile(jwtSigningKeyPath)
	if err != nil {
		logger.Error("read key failed", zap.Error(err))
		panic(err)
	}

	block, _ := pem.Decode(privateKey)

	key, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		logger.Error("private key parse failed", zap.Error(err))
		panic(err)
	}

	return &Secret{
		privateKey:      key,
		publicKey:       &key.PublicKey,
		jwtIssuer:       jwtIssuer,
		signingAlgoName: jwt.SigningMethodES512.Name,
		logger:          logger,
	}
}

func (r *Secret) CreateToken(username string) (*entity.Token, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES512,
		jwt.MapClaims{
			"jti": uuid.NewString(),
			"iss": r.jwtIssuer,
			"exp": time.Now().Add(defaultAccessExpiration).Unix(),
			"uid": username,
		},
	)

	accessToken, err := token.SignedString(r.PrivateKey())
	if err != nil {
		return nil, fmt.Errorf("jwt:s %w", err)
	}

	return pkg.PointerTo[entity.Token](entity.Token(accessToken)), nil
}

func (r *Secret) ParseJWT(token string) (jwt.MapClaims, error) {
	var mapClaims jwt.MapClaims

	_, err := jwt.ParseWithClaims(token, &mapClaims, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != r.signingAlgoName {
			return nil, ErrInvalidSigningAlgo
		}

		return r.PublicKey(), nil
	})
	if err != nil {
		return nil, err
	}

	return mapClaims, nil

}
