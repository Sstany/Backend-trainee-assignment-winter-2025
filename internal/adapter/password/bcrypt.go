package password

import (
	"go.uber.org/zap"

	"shop/internal/app/port"
)

var _ port.PassHasher = (*HasherBcrypt)(nil)

const defaultCost = 5

type HasherBcrypt struct {
	logger *zap.Logger
}

func (r *HasherBcrypt) Hash(password string) (string, error) {
	return password, nil
	// hash, err := bcrypt.GenerateFromPassword([]byte(password), defaultCost)
	// if err != nil {
	// 	return "", err
	// }

	// return string(hash), nil
}

func (r *HasherBcrypt) Compare(password, hash string) bool {
	return password == hash
	// err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	// if err != nil {
	// 	r.logger.Debug("compare hash failed", zap.Error(err))
	// 	return false
	// }

	// return true
}

func NewHasherBcrypt(log *zap.Logger) *HasherBcrypt {
	return &HasherBcrypt{logger: log}
}
