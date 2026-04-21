package shared

import "errors"

var ErrEmptyValue = errors.New("empty value error")

type ErrorCode int

const (
	ErrInvalidInput ErrorCode = 400
	ErrNotFound     ErrorCode = 404
	ErrUnauthorized ErrorCode = 401
	ErrForbidden    ErrorCode = 403
	ErrInternal     ErrorCode = 500
)

type Error struct {
	Code    ErrorCode
	Message string
	Err     error
}

func NewError(code ErrorCode, message string, err error) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func (e Error) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}

	return e.Message
}
