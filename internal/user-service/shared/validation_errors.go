package shared

type ValidationError struct {
	// Messages maps field names to a list of human-readable error messages.
	Messages map[string][]string `json:"messages"`
}

// Error implements the error interface for ValidationError.
func (e ValidationError) Error() string {
	return ServiceErrorCodeValidationError
}

const (
	// ServiceErrorCodeInternalError represents an internal error occurred during request handling
	ServiceErrorCodeInternalError string = "INTERNAL_ERROR"

	// ServiceErrorCodeNotFound not found.
	ServiceErrorCodeNotFound string = "NOT_FOUND"

	// ServiceErrorCodeInvalidFormat represents an invalid request body.
	ServiceErrorCodeInvalidFormat string = "INVALID_FORMAT"

	// ServiceErrorCodeInvalidRequest represents an invalid request error.
	ServiceErrorCodeInvalidRequest string = "INVALID_REQUEST"

	// ServiceErrorCodeValidationError represents a validation error.
	ServiceErrorCodeValidationError string = "VALIDATION_ERROR"
)
