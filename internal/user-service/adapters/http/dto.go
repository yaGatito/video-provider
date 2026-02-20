package httpadp

import (
	"strings"
	"video-provider/internal/user-service/shared"
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
	emailLen := len(r.Email)
	nameLen := len(r.Name)
	lastNameLen := len(r.LastName)
	passwordLen := len(r.Password)

	msgs := map[string][]string{}

	if emailLen == 0 {
		msgs["email"] = append(msgs["email"], "Email не може бути порожнім")
	}
	if emailLen > 100 {
		msgs["email"] = append(msgs["email"], "Email не може бути довшим за 100 символів")
	}
	if nameLen == 0 {
		msgs["name"] = append(msgs["name"], "Ім'я не може бути порожнім")
	}
	if nameLen > 50 {
		msgs["name"] = append(msgs["name"], "Ім'я не може бути довшим за 50 символів")
	}
	if lastNameLen == 0 {
		msgs["lastname"] = append(msgs["lastname"], "Прізвище не може бути порожнім")
	}
	if lastNameLen > 100 {
		msgs["lastname"] = append(msgs["lastname"], "Прізвище не може бути довшим за 100 символів")
	}
	if passwordLen < 8 {
		msgs["password"] = append(msgs["password"], "Пароль має бути не менше 8 символів")
	}
	if passwordLen > 100 {
		msgs["password"] = append(msgs["password"], "Пароль не може бути довшим за 100 символів")
	}
	if len(msgs) > 0 {
		return shared.ValidationError{Messages: msgs}
	}
	return nil
}

func (r createUserRequest) normalize() {
	r.Email = strings.TrimSpace(strings.ToLower(r.Email))
	r.Name = strings.TrimSpace(r.Name)
	r.LastName = strings.TrimSpace(r.LastName)
	r.Password = strings.TrimSpace(r.Password)
}

// serviceErrorResponse represents a generic error response from the service.
type serviceErrorResponse struct {
	Code    string `json:"code"`
	Payload any    `json:"payload,omitempty"`
}
