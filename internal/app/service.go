package app

import (
	"fmt"
	"log"
	"time"
	"user-service-DDD/internal/domain"
	"user-service-DDD/internal/ports"
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

type UserService interface {
	Register(cmd RegisterUserCommand) (int64, error)
	Get(id int64) (GetUserResult, error)
}

type UserUsecasesManager struct {
	Repo ports.UserRepository
}

func NewUserUsecasesManager(repo ports.UserRepository) UserUsecasesManager {
	return UserUsecasesManager{Repo: repo}
}

func (ui UserUsecasesManager) Register(cmd RegisterUserCommand) (int64, error) {
	user, err := domain.NewUser(cmd.Email, cmd.Name, cmd.Lastname)
	if err != nil {
		return -1, err
	}
	fmt.Println("Received RegisterUserCommand with valid email", user.Email)

	pass := domain.Password(cmd.Password)
	if err = pass.Validate(); err != nil {
		return -1, err
	}

	// TODO: use caching mechanism
	var passHash = string(pass)
	var passSalt = string(pass)

	id, err := ui.Repo.Create(user, passHash, passSalt)
	if err != nil {
		log.Printf("Error creating user: %v\n", err)
		return 0, err
	}

	log.Printf("User created and saved into DB with ID: %d\n", id)
	return id, nil
}

func (ui UserUsecasesManager) Get(id int64) (GetUserResult, error) {
	user, err := ui.Repo.FindByID(id)
	if err != nil {
		log.Printf("Error retrieving user with ID %d: %v\n", id, err)
		return GetUserResult{}, err
	}
	return GetUserResult{
		Email:    user.Email,
		Name:     user.Name,
		Lastname: user.LastName,
		CreateAt: user.CreatedAt,
	}, nil
}
