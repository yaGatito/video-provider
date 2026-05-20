package app

import (
	"context"
	"video-provider/pkg/auth"
	"video-provider/pkg/common"
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
	if user == domain.Nil {
		return uuid.Nil, &common.Error{
			Code:    common.ErrInvalidInput,
			Message: "empty user"}
	}
	if password == "" {
		return uuid.Nil, &common.Error{
			Code: common.ErrInvalidInput, Message: "empty password"}
	}

	hash, err := us.Hasher.Hash(password)
	if err != nil {
		return uuid.Nil, &common.Error{
			Err:     err,
			Code:    common.ErrInternal,
			Message: "failed to hash password"}
	}

	id, err := us.Repo.Create(ctx, user, hash)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (us *UserService) Get(ctx context.Context, id uuid.UUID) (domain.User, error) {
	if id == uuid.Nil {
		return domain.Nil, &common.Error{
			Code:    common.ErrInvalidInput,
			Message: "empty id"}
	}
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
		return "", &common.Error{
			Code: common.ErrInvalidInput, Message: "empty email"}
	}
	if len(password) == 0 {
		return "", &common.Error{
			Code:    common.ErrInvalidInput,
			Message: "empty password"}
	}

	userID, hash, err := us.Repo.GetPasswordHash(ctx, email)
	if err != nil {
		return "", err
	}

	err = us.Hasher.CompareHashAndPassword(hash, password)
	if err != nil {
		return "", &common.Error{
			Err:  err,
			Code: common.ErrUnauthorized, Message: "failed to compare password"}
	}

	token, err := us.Tokenizer.CreateToken(userID)
	if err != nil {
		return "", err
	}

	return token, nil
}
