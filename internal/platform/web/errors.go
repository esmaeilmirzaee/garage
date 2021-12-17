package web

import "github.com/pkg/errors"

// FieldError is used to indicate an error with a specific request field.
type FieldError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

// ErrorResponse how we respond to clients when something goes wrong
// ErrorResponse is the form used for API responses from failures in the API.
type ErrorResponse struct {
	Error  string       `json:"error"`
	Fields []FieldError `json:"fields,omitempty"`
}

// Error is used to pass an error during the request through the
// application with web specific request.
type Error struct {
	Err    error
	Status int
	Fields []FieldError
}

// NewRequestError is used when a known error condition is encountered.
func NewRequestError(err error, status int) *Error {
	return &Error{
		Err:    err,
		Status: status,
	}
}

// Error return string type of error message
func (e *Error) Error() string {
	return e.Err.Error()
}

// Shutdown is a type used to help with the graceful termination of
// the service
type shutdown struct {
	Message string
}

// Error is the implementation of the Error interface.
func (s *shutdown) Error() string {
	return s.Message
}

// NewShutdownError returns an error that causes the framework to
// signal a graceful shutdown.
func NewShutdownError(message string) error {
	return &shutdown{message}
}

// IsShutdown checks to see if the shutdown error is contained in
// the specified error value.
func IsShutdown(err error) bool {
	if _, ok := errors.Cause(err).(*shutdown); ok {
		return true
	}
	return false
}
