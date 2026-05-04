package common

import (
	"errors"
	"fmt"
)

var ErrEmptyValue = errors.New("empty value error")

type ErrorCode int

// TODO: reconsider internal errors codes to be able distinguish em
const (
	ErrInvalidInput ErrorCode = 400
	ErrUnauthorized ErrorCode = 401
	ErrForbidden    ErrorCode = 403
	ErrNotFound     ErrorCode = 404
	ErrInternal     ErrorCode = 500
)

type Error struct {
	Code    ErrorCode
	Message string
	Details any
	Err     error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}

	return e.Message
}

func (e Error) GetDetails() string {
	if e.Details != nil {
		return fmt.Sprintf("%+v", e.Details)
	}

	return ""
}
