package app

import (
	"fmt"
	"log"
	"time"
	"video-provider/internal/user-service/domain"
	"video-provider/internal/user-service/ports"

	"github.com/google/uuid"
)

type RegisterUserCommand struct {
	Email    string
	Name     string
	Lastname string
	Password string
}

type GetUserResult struct {
	Name     string
	Email    string
	Lastname string
	CreateAt time.Time
}

type UserInteractor interface {
	Register(cmd RegisterUserCommand) (uuid.UUID, error)
	Get(id uuid.UUID) (GetUserResult, error)
}

type UserService struct {
	Repo ports.UserRepository
	log  log.Logger
}

func NewUserService(repo ports.UserRepository) *UserService {
	return &UserService{Repo: repo}
}

func (us *UserService) Register(cmd RegisterUserCommand) (uuid.UUID, error) {
	user, err := domain.NewUser(cmd.Email, cmd.Name, cmd.Lastname)
	if err != nil {
		return uuid.UUID{}, err
	}
	fmt.Println("Received RegisterUserCommand with valid email", user.Email)

	pass := domain.Password(cmd.Password)
	if err = pass.ValidatePassword(); err != nil {
		return uuid.UUID{}, fmt.Errorf("error validating user: %w", err)
	}

	// TODO: use hashing mechanism
	var passHash = string(pass)
	var passSalt = string(pass)

	id, err := us.Repo.Create(user, passHash, passSalt)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("error creating user: %w", err)
	}

	log.Printf("User created and saved into DB with ID: %s\n", id.String())
	return id, nil
}

func (us *UserService) Get(id uuid.UUID) (GetUserResult, error) {
	user, err := us.Repo.FindByID(id)
	if err != nil {
		log.Printf("Error retrieving user with ID %s: %v\n", id.String(), err)
		return GetUserResult{}, err
	}
	return GetUserResult{
		Email:    user.Email,
		Name:     user.Name,
		Lastname: user.LastName,
		CreateAt: user.CreatedAt,
	}, nil
}
