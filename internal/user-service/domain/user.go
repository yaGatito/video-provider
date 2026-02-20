package domain

import (
	"regexp"
	"time"
	"video-provider/internal/user-service/shared"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	Email     string
	Name      string
	LastName  string
	CreatedAt time.Time // 2025-08-14 00:37:00
	IsAdmin   bool      // internal-only; guarded by service logic
	Status    string    // "active", "disabled"
}

func NewUser(email string, name string, lastname string) (User, error) {
	// Regex validation
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	nameRegex := regexp.MustCompile(`^[a-zA-Z]+$`)     // only letters
	lastNameRegex := regexp.MustCompile(`^[a-zA-Z]+$`) // only letters

	if !emailRegex.MatchString(email) {
		return User{}, shared.ServiceError{
			Code: shared.InvalidFormatErr,
			Msg:  "invalid email format"}
	}
	if !nameRegex.MatchString(name) {
		return User{}, shared.ServiceError{
			Code: shared.InvalidFormatErr,
			Msg:  "invalid name format",
		}
	}
	if !lastNameRegex.MatchString(lastname) {
		return User{}, shared.ServiceError{
			Code: shared.InvalidFormatErr,
			Msg:  "invalid lastname format"}
	}

	return User{
		Email:     email,
		Name:      name,
		LastName:  lastname,
		CreatedAt: time.Now(),
		IsAdmin:   false,
		Status:    "active",
	}, nil
}

type Password string

func (p Password) ValidatePassword() error {
	passRegex := regexp.MustCompile(`^[a-zA-Z0-9]{8,255}$`)

	matchString := passRegex.MatchString(string(p))
	if matchString {
		return nil
	} else {
		return shared.ServiceError{Code: shared.InvalidFormatErr, Msg: "password must be 8 characters long and contain at least one uppercase letter and "}
	}
}
