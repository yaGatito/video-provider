package ports

import (
	"examples-user-service/internal/domain"
)

type UserRepository interface {
	Create(user *domain.User, passwordHash string, passwordSalt string) (int64, error)
	FindByID(id int64) (*domain.User, error)
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
