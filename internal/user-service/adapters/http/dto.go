package httpadp

import (
	"strings"
	"video-provider/internal/user-service/shared"

	"github.com/go-playground/validator/v10"
)

// createUserRequest represents the data needed to create a new user.
// separating this struct in order to abstract service from transport (json tags required)
type createUserRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	LastName string `json:"lastname"`
	Password string `json:"password"`
}

// validate createUserRequest validates the createUserRequest fields.
// It checks for empty fields, length constraints, and returns a validationError
func (r createUserRequest) validate() error {
	if len(r.Email) == 0 {
		return shared.ServiceError{Code: shared.InvalidFormatErr, Msg: "Email shouldn't be empty"}
	}
	if len(r.Email) > 100 {
		return shared.ServiceError{Code: shared.InvalidFormatErr, Msg: "Email too long"}
	}

	if len(r.Name) == 0 {
		return shared.ServiceError{Code: shared.InvalidFormatErr, Msg: "Name shouldn't be empty"}
	}
	if len(r.Name) > 50 {
		return shared.ServiceError{Code: shared.InvalidFormatErr, Msg: "Name too long"}
	}

	if len(r.LastName) == 0 {
		return shared.ServiceError{Code: shared.InvalidFormatErr, Msg: "Lastname shouldn't be empty"}
	}
	if len(r.LastName) > 100 {
		return shared.ServiceError{Code: shared.InvalidFormatErr, Msg: "Lastname too long"}
	}

	if len(r.Password) < 8 {
		return shared.ServiceError{Code: shared.InvalidFormatErr, Msg: "Password too short"}
	}
	if len(r.Password) > 100 {
		return shared.ServiceError{Code: shared.InvalidFormatErr, Msg: "Password too long"}
	}

	return nil
}

// normalize normalizes the user request data.
func (r createUserRequest) normalize() {
	r.Email = strings.TrimSpace(strings.ToLower(r.Email))
	r.Name = strings.TrimSpace(r.Name)
	r.LastName = strings.TrimSpace(r.LastName)
	// TODO: check spaces in password and trim them if present.
	// TODO: refactor validation.
	// r.Password = strings.TrimSpace(r.Password)
}

// serviceErrorResponse represents a generic error response from the service.
type serviceErrorResponse struct {
	Code    string `json:"code"`
	Payload any    `json:"payload,omitempty"`
}

// loginUserRequest represents the request body for login
type loginUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=50"`
}

// authResponse represents the response with authentication token
type authResponse struct {
	Token string `json:"token"`
}

// validate checks if the login request is valid
func (r *loginUserRequest) validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	err := validate.Struct(r)
	if err != nil {
		// Handle validation errors
		return err
	}
	return nil
}

// normalize normalizes the login request
func (r *loginUserRequest) normalize() {
	r.Email = strings.ToLower(r.Email)
}
