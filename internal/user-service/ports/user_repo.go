package ports

import (
	"video-provider/internal/user-service/domain"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(user domain.User, passwordHash string, passwordSalt string) (uuid.UUID, error)
	Update(user domain.User) error
	FindByID(id uuid.UUID) (domain.User, error)
	FindByEmail(email string) (domain.User, error)
}

type PasswordHasher interface {
	Hash(password string) ([]byte, error)
}
