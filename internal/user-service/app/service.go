package app

import (
	"context"
	"log"
	"video-provider/common/auth"
	"video-provider/common/shared"
	"video-provider/user-service/domain"
	"video-provider/user-service/ports"

	"github.com/google/uuid"
)

type UserInteractor interface {
	Create(ctx context.Context, user domain.User, password string) (uuid.UUID, error)
	Get(ctx context.Context, id uuid.UUID) (domain.User, error)
	Update(ctx context.Context, id uuid.UUID, user domain.User) error
	Login(ctx context.Context, email string, password []byte) (string, error)
}

type UserService struct {
	Repo      ports.UserRepository
	Hasher    ports.PasswordHasher
	log       log.Logger
	Tokenizer *auth.Tokenizer
}

func NewUserService(
	repo ports.UserRepository,
	hasher ports.PasswordHasher,
	tokenizer *auth.Tokenizer,
) *UserService {
	return &UserService{Repo: repo, Hasher: hasher, Tokenizer: tokenizer}
}

func (us *UserService) Create(
	ctx context.Context,
	user domain.User,
	password string,
) (uuid.UUID, error) {
	hash, err := us.Hasher.Hash(password)
	if err != nil {
		return uuid.UUID{}, shared.NewError(shared.ErrInternal, "failed to hash password", err)
	}

	id, err := us.Repo.Create(ctx, user, hash)
	if err != nil {
		return uuid.UUID{}, err
	}

	return id, nil
}

func (us *UserService) Get(ctx context.Context, id uuid.UUID) (domain.User, error) {
	return us.Repo.FindByID(ctx, id)
}

func (us *UserService) Update(ctx context.Context, id uuid.UUID, toUpdate domain.User) error {
	user, err := us.Repo.FindByID(ctx, id)
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

	err = us.Repo.Update(ctx, id, user)
	if err != nil {
		return err
	}

	return nil
}

func (us *UserService) Login(ctx context.Context, email string, password []byte) (string, error) {
	if email == "" {
		return "", shared.NewError(shared.ErrInvalidInput, "email is required", nil)
	}
	if len(password) == 0 {
		return "", shared.NewError(shared.ErrInvalidInput, "password is required", nil)
	}

	userID, hash, err := us.Repo.GetPasswordHash(ctx, email)
	if err != nil {
		return "", err
	}

	err = us.Hasher.CompareHashAndPassword(hash, password)
	if err != nil {
		return "", shared.NewError(shared.ErrUnauthorized, "failed to compare password", nil)
	}

	token, err := us.Tokenizer.CreateToken(userID)
	if err != nil {
		return "", err
	}

	return token, nil
}
