package port

import "crypto/ecdsa"

type SecretRepo interface {
	PublicKey() *ecdsa.PublicKey
	PrivateKey() *ecdsa.PrivateKey
	JWTIssuer() string
}
