package app

import (
	"fmt"
	"log"
	"time"
	"video-provider/internal/user-service/domain"
	"video-provider/internal/user-service/ports"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type RegisterUserCommand struct {
	Email    string
	Name     string
	Lastname string
	Password string
}

// LoginUserCommand represents the data needed to login a user.
type LoginUserCommand struct {
	Email    string
	Password string
}

type UpdateUserCommand struct {
	ID       uuid.UUID
	Email    *string
	Name     *string
	Lastname *string
	Password *string
}

type GetUserResult struct {
	Name     string
	Email    string
	Lastname string
	CreateAt time.Time
}

type UserInteractor interface {
	Create(cmd RegisterUserCommand) (uuid.UUID, error)
	Get(id uuid.UUID) (GetUserResult, error)
	Update(cmd UpdateUserCommand) error
	Login(email, password string) (string, error)
}

type UserService struct {
	Repo ports.UserRepository
	log  log.Logger
}

func NewUserService(repo ports.UserRepository) *UserService {
	return &UserService{Repo: repo}
}

func (us *UserService) Create(cmd RegisterUserCommand) (uuid.UUID, error) {
	user := domain.NewUser(cmd.Email, cmd.Name, cmd.Lastname)
	fmt.Println("Received RegisterUserCommand with valid email", user.Email)

	// TODO: use hashing mechanism
	id, err := us.Repo.Create(user, cmd.Password, cmd.Password)
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

func (us *UserService) Update(cmd UpdateUserCommand) error {
	user, err := us.Repo.FindByID(cmd.ID)
	if err != nil {
		return fmt.Errorf("error finding user by ID: %w", err)
	}

	if cmd.Email != nil {
		user.Email = *cmd.Email
	}
	if cmd.Name != nil {
		user.Name = *cmd.Name
	}
	if cmd.Lastname != nil {
		user.LastName = *cmd.Lastname
	}
	if cmd.Password != nil {
		// pass := domain.Password(*cmd.Password)
		// if err := pass.ValidatePassword(); err != nil {
		// 	return fmt.Errorf("error validating user password: %w", err)
		// }
		// var passHash = string(pass)
		// var passSalt = string(pass)

		// err := us.Repo.UpdatePass(cmd.ID, passHash, passSalt)
		// if err != nil {
		// 	return fmt.Errorf("error updating user password: %w", err)
		// }
	}

	err = us.Repo.Update(user)
	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	log.Printf("User updated with ID: %s\n", cmd.ID.String())
	return nil
}

func (us *UserService) Login(email, password string) (string, error) {
	// Find user by email
	user, err := us.Repo.FindByEmail(email)
	if err != nil {
		log.Printf("Error finding user by email: %v\n", err)
		return "", fmt.Errorf("error finding user by email: %w", err)
	}

	// Validate password
	// if !user.PasswordMatches(password) {
	// 	return "", fmt.Errorf("invalid credentials")
	// }

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})

	// Sign token with a secret key (you should use a secure key in production)
	secretKey := "your-secret-key-here"
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("error generating token: %w", err)
	}

	return tokenString, nil
}
