package shared

import "errors"

var ErrEmptyValue = errors.New("empty value error")

type ServiceError struct {
	Code    int
	Message string
	Err     error
}

// Error implements the error interface for InvalidRequestError.
func (e ServiceError) Error() string {
	return e.Message
}
