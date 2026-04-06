package cryptoadp

import (
	"video-provider/internal/user-service/ports"

	"golang.org/x/crypto/bcrypt"
)

type BcryptPasswordHasher struct{}

// Ensure PostgresUserRepository implements ports.UserRepository
var _ ports.PasswordHasher = (*BcryptPasswordHasher)(nil)

func NewBCryptPasswordHasher() *BcryptPasswordHasher {
	return &BcryptPasswordHasher{}
}

func (h *BcryptPasswordHasher) Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func (h *BcryptPasswordHasher) CompareHashAndPassword(hash, password []byte) error {
	return bcrypt.CompareHashAndPassword(hash, password)
}
