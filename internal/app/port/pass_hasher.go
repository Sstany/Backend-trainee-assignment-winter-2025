package port

//go:generate mockgen -destination ../../adapter/password/hasher_mock.go -package password -source ./pass_hasher.go

type PassHasher interface {
	Hash(password string) (string, error)
	Compare(password, hash string) bool
}
