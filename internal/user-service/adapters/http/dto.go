package httpadp

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

// createUserRequest represents the data needed to create a new user.
// separating this struct in order to abstract service from transport (json tags required)
type createUserRequest struct {
	Email    string `json:"email" validate:"required,email,lenLimit"`
	Name     string `json:"name" validate:"required,text64"`
	LastName string `json:"lastname" validate:"required,text64"`
	Password string `json:"password" validate:"required"`
}

// normalize normalizes the user request data.
func (r createUserRequest) normalize() {
	r.Email = strings.TrimSpace(strings.ToLower(r.Email))
	r.Name = strings.TrimSpace(r.Name)
	r.LastName = strings.TrimSpace(r.LastName)
}

// serviceErrorResponse represents a generic error response from the service.
type serviceErrorResponse struct {
	Code    string `json:"code"`
	Payload any    `json:"payload,omitempty"`
}

// loginUserRequest represents the request body for login
type loginUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
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
