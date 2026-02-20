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

func NewUser(email string, name string, lastname string) (*User, error) {
	msgs := map[string][]string{}

	// Regex validation
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	nameRegex := regexp.MustCompile(`^[a-zA-Z]+$`)     // only letters
	lastNameRegex := regexp.MustCompile(`^[a-zA-Z]+$`) // only letters

	if !emailRegex.MatchString(email) {
		msgs["email"] = append(msgs["email"], "invalid_format")
	}
	if !nameRegex.MatchString(name) {
		msgs["name"] = append(msgs["name"], "invalid_format")
	}
	if !lastNameRegex.MatchString(lastname) {
		msgs["lastname"] = append(msgs["lastname"], "invalid_format")
	}

	if len(msgs) > 0 {
		return nil, shared.ValidationError{Messages: msgs}
	}

	return &User{
		Email:     email,
		Name:      name,
		LastName:  lastname,
		CreatedAt: time.Now(),
		IsAdmin:   false,
		Status:    "active",
	}, nil
}

type Password string

func (p Password) Validate() error {
	passRegex := regexp.MustCompile(`^[a-zA-Z0-9]{8,255}$`)

	matchString := passRegex.MatchString(string(p))
	if matchString {
		return nil
	} else {
		return shared.ValidationError{Messages: map[string][]string{"password": {"invalid_format"}}}
	}
}
