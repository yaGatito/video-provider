package app

import (
	"context"
	"log"
	"video-provider/internal/pkg/shared"
	"video-provider/internal/user-service/domain"
	"video-provider/internal/user-service/ports"

	"github.com/google/uuid"
)

type UserInteractor interface {
	Create(ctx context.Context, user domain.User, password string) (uuid.UUID, error)
	Get(ctx context.Context, id uuid.UUID) (domain.User, error)
	Update(ctx context.Context, id uuid.UUID, user domain.User) error
	Login(ctx context.Context, email, password string) (string, error)
}

type UserService struct {
	repo   ports.UserRepository
	hasher ports.PasswordHasher
	log    log.Logger
}

func NewUserService(repo ports.UserRepository, hasher ports.PasswordHasher) *UserService {
	return &UserService{repo: repo, hasher: hasher}
}

func (us *UserService) Create(ctx context.Context, user domain.User, password string) (uuid.UUID, error) {
	hash, err := us.hasher.Hash(password)
	if err != nil {
		return uuid.UUID{}, shared.NewError(shared.ErrInternal, "failed to hash password", err)
	}

	id, err := us.repo.Create(ctx, user, hash)
	if err != nil {
		return uuid.UUID{}, err
	}

	return id, nil
}

func (us *UserService) Get(ctx context.Context, id uuid.UUID) (domain.User, error) {
	user, err := us.repo.FindByID(ctx, id)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (us *UserService) Update(ctx context.Context, id uuid.UUID, toUpdate domain.User) error {
	user, err := us.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if toUpdate.Email != "" {
		user.Email = toUpdate.Email
	}
	if toUpdate.Name != "" {
		user.Name = toUpdate.Name
	}
	if toUpdate.LastName != "" {
		user.LastName = toUpdate.LastName
	}

	err = us.repo.Update(ctx, id, user)
	if err != nil {
		return err
	}

	return nil
}

func (us *UserService) Login(ctx context.Context, email, password string) (string, error) {
	_, hash, err := us.repo.GetPassword(ctx, email)
	if err != nil {
		return "", err
	}

	err = us.hasher.CompareHashAndPassword(hash, password)
	if err != nil {
		return "", shared.NewError(shared.ErrUnauthorized, "failed to compare password", err)
	}

	// TODO: Generate and return a JWT token or session ID here

	return "success", nil
}
