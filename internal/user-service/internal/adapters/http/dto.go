package httpadp

import (
	"strings"
	"user-service/internal/shared"
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
	var v []shared.FieldViolationError

	emailLen := len(r.Email)
	nameLen := len(r.Name)
	lastNameLen := len(r.LastName)
	passwordLen := len(r.Password)

	if emailLen == 0 {
		v = append(v, shared.FieldViolationError{
			ViolatedField: shared.ViolatedFieldEmail, ViolationCode: shared.ViolationCodeEmpty, Message: "Email не може бути порожнім",
		})
	}
	if emailLen > 100 {
		v = append(v, shared.FieldViolationError{
			ViolatedField: shared.ViolatedFieldEmail, ViolationCode: shared.ViolationCodeTooLong, Message: "Email не може бути довшим за 100 символів",
		})
	}
	if nameLen == 0 {
		v = append(v, shared.FieldViolationError{
			ViolatedField: shared.ViolatedFieldName, ViolationCode: shared.ViolationCodeEmpty, Message: "Ім'я не може бути порожнім",
		})
	}
	if nameLen > 50 {
		v = append(v, shared.FieldViolationError{
			ViolatedField: shared.ViolatedFieldName, ViolationCode: shared.ViolationCodeTooLong, Message: "Ім'я не може бути довшим за 50 символів",
		})
	}
	if lastNameLen == 0 {
		v = append(v, shared.FieldViolationError{
			ViolatedField: shared.ViolatedFieldLastName, ViolationCode: shared.ViolationCodeEmpty, Message: "Прізвище не може бути порожнім",
		})
	}
	if lastNameLen > 100 {
		v = append(v, shared.FieldViolationError{
			ViolatedField: shared.ViolatedFieldLastName, ViolationCode: shared.ViolationCodeTooLong, Message: "Прізвище не може бути довшим за 100 символів",
		})
	}
	if passwordLen < 8 {
		v = append(v, shared.FieldViolationError{
			ViolatedField: shared.ViolatedFieldPassword, ViolationCode: shared.ViolationCodeTooShort, Message: "Пароль має бути не менше 8 символів",
		})
	}
	if passwordLen > 100 {
		v = append(v, shared.FieldViolationError{
			ViolatedField: shared.ViolatedFieldPassword, ViolationCode: shared.ViolationCodeTooLong, Message: "Пароль не може бути довшим за 100 символів",
		})
	}
	if len(v) > 0 {
		return shared.ValidationError{
			Violations: v,
		}
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
