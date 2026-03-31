package httpadp

import (
	"fmt"
	"net/http"
	"regexp"
	"video-provider/internal/pkg/shared"
	"video-provider/internal/user-service/policy"

	"github.com/go-playground/validator/v10"
)

var textRe = regexp.MustCompile(fmt.Sprintf(`^[A-Za-z]{%d,%d}$`, policy.MinInputTextLen, policy.MaxInputTextLen))

func NewUserValidate() *validator.Validate {
	validate := validator.New(validator.WithRequiredStructEnabled())

	validate.RegisterValidation("text64", func(fl validator.FieldLevel) bool {
		return textRe.MatchString(fl.Field().String())
	})
	validate.RegisterValidation("lenLimit", func(fl validator.FieldLevel) bool {
		return len(fl.Field().String()) <= policy.MaxInputTextLen
	})

	return validate
}

var passRe = regexp.MustCompile(`^[A-Za-z0-9!@#$%^&*()_=+]+$`)
var passReqDigGroup = regexp.MustCompile(`.*[0-9]+.*`)
var passReqLowGroup = regexp.MustCompile(`.*[a-z]+.*`)
var passReqCapGroup = regexp.MustCompile(`.*[A-Z]+.*`)
var passReqSpecGroup = regexp.MustCompile(`.*[!@#$%^&*()_=+-]+.*`)

func validatePassword(password []byte) error {
	if len(password) < policy.MinPasswordLen || len(password) > policy.MaxInputTextLen {
		return shared.ServiceError{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("password must be between %d and %d characters long", policy.MinPasswordLen, policy.MaxInputTextLen),
		}
	}
	if !passRe.Match(password) {
		return shared.ServiceError{
			Code:    http.StatusBadRequest,
			Message: "password contains invalid characters",
		}
	}
	if !passReqDigGroup.Match(password) {
		return shared.ServiceError{
			Code:    http.StatusBadRequest,
			Message: "password must contain at least one digit",
		}
	}
	if !passReqLowGroup.Match(password) {
		return shared.ServiceError{
			Code:    http.StatusBadRequest,
			Message: "password must contain at least one lowercase letter",
		}
	}
	if !passReqCapGroup.Match(password) {
		return shared.ServiceError{
			Code:    http.StatusBadRequest,
			Message: "password must contain at least one uppercase letter",
		}
	}
	if !passReqSpecGroup.Match(password) {
		return shared.ServiceError{
			Code:    http.StatusBadRequest,
			Message: "password must contain at least one special character",
		}
	}
	return nil
}
