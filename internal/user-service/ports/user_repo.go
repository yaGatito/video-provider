package ports

import (
	"video-provider/internal/user-service/domain"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(user *domain.User, passwordHash string, passwordSalt string) (uuid.UUID, error)
	FindByID(id uuid.UUID) (*domain.User, error)
}

type PasswordHasher interface {
	Hash(password string) ([]byte, error)
}

// TODO: implement
//type IDGen interface {
//	NewUserID() domain.UserID
//}

//type Clock interface {
//	Now() int64 // або time.Time; оберемо time.Time на app-рівні нижче
//}
