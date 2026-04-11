package app

import (
	"context"
	"log"
	"time"
	"user-service/domain"
	"user-service/pkg/shared"
	"user-service/ports"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type UserInteractor interface {
	Create(ctx context.Context, user domain.User, password string) (uuid.UUID, error)
	Get(ctx context.Context, id uuid.UUID) (domain.User, error)
	Update(ctx context.Context, id uuid.UUID, user domain.User) error
	Login(ctx context.Context, email string, password []byte) (string, error)
}

type UserService struct {
	repo         ports.UserRepository
	hasher       ports.PasswordHasher
	log          log.Logger
	getJWTSecret func() []byte
}

func NewUserService(
	repo ports.UserRepository,
	hasher ports.PasswordHasher,
	getSecret func() []byte,
) *UserService {
	return &UserService{repo: repo, hasher: hasher, getJWTSecret: getSecret}
}

func (us *UserService) Create(
	ctx context.Context,
	user domain.User,
	password string,
) (uuid.UUID, error) {
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
	return us.repo.FindByID(ctx, id)
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

func (us *UserService) Login(ctx context.Context, email string, password []byte) (string, error) {
	if email == "" {
		return "", shared.NewError(shared.ErrInvalidInput, "email is required", nil)
	}
	if len(password) == 0 {
		return "", shared.NewError(shared.ErrInvalidInput, "password is required", nil)
	}

	userID, hash, err := us.repo.GetPasswordHash(ctx, email)
	if err != nil {
		return "", err
	}

	err = us.hasher.CompareHashAndPassword(hash, password)
	if err != nil {
		return "", shared.NewError(shared.ErrUnauthorized, "failed to compare password", nil)
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(us.getJWTSecret())

	if err != nil {
		return "", err
	}

	return signedToken, nil
}
