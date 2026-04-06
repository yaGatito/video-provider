package ports

import (
	"context"
	"video-provider/internal/user-service/domain"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user domain.User, password []byte) (uuid.UUID, error)
	Update(ctx context.Context, id uuid.UUID, user domain.User) error
	FindByID(ctx context.Context, id uuid.UUID) (domain.User, error)
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	GetPasswordHash(ctx context.Context, email string) (uuid.UUID, []byte, error)
}
