package repo

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"os"

	"shop/internal/app/port"

	"go.uber.org/zap"
)

var _ port.SecretRepo = (*Secret)(nil)

type Secret struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
	jwtIssuer  string
	logger     *zap.Logger
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
		privateKey: key,
		publicKey:  &key.PublicKey,
		jwtIssuer:  jwtIssuer,
		logger:     logger,
	}
}
