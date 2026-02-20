package shared

const (
	// InternalErr represents an internal error occurred during request handling
	InternalErr string = "INTERNAL_ERROR"

	// NotFoundErr not found.
	NotFoundErr string = "NOT_FOUND"

	// InvalidFormatErr represents an invalid request body.
	InvalidFormatErr string = "INVALID_FORMAT"

	// InvalidRequestErr represents an invalid request error.
	InvalidRequestErr string = "INVALID_REQUEST"

	// ValidationErr represents a validation error.
	ValidationErr string = "VALIDATION_ERROR"
)

type ServiceError struct {
	Code string
	Msg  string
}

// Error implements the error interface for ValidationError.
func (e ServiceError) Error() string {
	return e.Code
}
