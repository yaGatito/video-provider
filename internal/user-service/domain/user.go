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
	var valError = shared.ValidationError{
		//	Code:       shared.ServiceErrorCodeValidationError,
		Violations: []shared.FieldViolationError{},
	}

	// Regex validation
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	nameRegex := regexp.MustCompile(`^[a-zA-Z]+$`)     // only letters
	lastNameRegex := regexp.MustCompile(`^[a-zA-Z]+$`) // only letters

	if !emailRegex.MatchString(email) {
		valError.Violations = append(valError.Violations, shared.FieldViolationError{
			ViolatedField: shared.ViolatedFieldEmail,
			ViolationCode: shared.ViolationCodeInvalidFormat,
		})
	}
	if !nameRegex.MatchString(name) {
		valError.Violations = append(valError.Violations, shared.FieldViolationError{
			ViolatedField: shared.ViolatedFieldName,
			ViolationCode: shared.ServiceErrorCodeInvalidFormat,
		})
	}
	if !lastNameRegex.MatchString(lastname) {
		valError.Violations = append(valError.Violations, shared.FieldViolationError{
			ViolatedField: shared.ViolatedFieldLastName,
			ViolationCode: shared.ViolationCodeInvalidFormat,
		})
	}

	if len(valError.Violations) > 0 {
		return &User{
			Email:     email,
			Name:      name,
			LastName:  lastname,
			CreatedAt: time.Now(),
			IsAdmin:   false,
			Status:    "active",
		}, nil
	} else {
		return nil, valError
	}
}

type Password string

func (p Password) Validate() error {
	passRegex := regexp.MustCompile(`^[a-zA-Z0-9]{8,255}$`)

	matchString := passRegex.MatchString(string(p))
	if matchString {
		return nil
	} else {
		return shared.ValidationError{
			Violations: []shared.FieldViolationError{{ViolationCode: shared.ViolationCodeInvalidFormat, ViolatedField: shared.ViolatedFieldPassword}},
		}
	}
}
