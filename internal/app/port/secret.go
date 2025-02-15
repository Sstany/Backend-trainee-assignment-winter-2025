package port

import (
	"crypto/ecdsa"

	"shop/internal/app/entity"

	"github.com/golang-jwt/jwt"
)

//go:generate mockgen -destination ../../adapter/repo/mock/secret_mock.go -package repo -source ./secret.go

type SecretRepo interface {
	PublicKey() *ecdsa.PublicKey
	PrivateKey() *ecdsa.PrivateKey
	JWTIssuer() string
	ParseJWT(token string) (jwt.MapClaims, error)
	CreateToken(username string) (*entity.Token, error)
}
