package shared

type ValidationError struct {
	Violations []FieldViolationError `json:"violations"`
}

// Error implements the error interface for validationError.
func (e ValidationError) Error() string {
	return ServiceErrorCodeValidationError
}

type FieldViolationError struct {
	ViolatedField string `json:"field"`
	ViolationCode string `json:"code"`
	Message       string `json:"message"`
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

const (
	ViolatedFieldEmail    string = "email"
	ViolatedFieldName     string = "name"
	ViolatedFieldLastName string = "lastname"
	ViolatedFieldPassword string = "password"
)

const (
	ViolationCodeEmpty         string = "EMPTY"
	ViolationCodeTooShort      string = "TOO_SHORT"
	ViolationCodeTooLong       string = "TOO_LONG"
	ViolationCodeInvalidFormat string = "INVALID_FORMAT"
)
