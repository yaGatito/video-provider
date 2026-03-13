package shared

const (
	// InternalErr represents an internal error occurred during request handling
	InternalErr ServiceCode = "INTERNAL_ERROR"

	// NotFoundErr not found.
	NotFoundErr ServiceCode = "NOT_FOUND"

	// InvalidFormatErr represents an invalid request body.
	InvalidFormatErr ServiceCode = "INVALID_FORMAT"

	// InvalidRequestErr represents an invalid request error.
	InvalidRequestErr ServiceCode = "INVALID_REQUEST"

	// UnauthorizedErr represents an authorization error.
	UnauthorizedErr ServiceCode = "UNAUTHORIZED_ERROR"
)

type ServiceCode string

type ServiceError struct {
	Code ServiceCode
	Msg  string
}

// Error implements the error interface for InvalidRequestError.
func (e ServiceError) Error() string {
	return e.Msg
}
