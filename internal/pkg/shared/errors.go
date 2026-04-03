package shared

import "errors"

var ErrEmptyValue = errors.New("empty value error")

type ErrorCode string

const (
	ErrInvalidInput ErrorCode = "ERR_INVALID_INPUT"
	ErrNotFound     ErrorCode = "ERR_NOT_FOUND"
	ErrInternal     ErrorCode = "ERR_INTERNAL"
	ErrUnauthorized ErrorCode = "ERR_UNAUTHORIZED"
)

type Error struct {
	Code    ErrorCode
	Message string
	Err     error
}

func NewError(code ErrorCode, message string, err error) Error {
	return Error{
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
