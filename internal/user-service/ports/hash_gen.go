package ports

type PasswordHasher interface {
	Hash(password string) ([]byte, error)
	CompareHashAndPassword(hash, password string) error
}