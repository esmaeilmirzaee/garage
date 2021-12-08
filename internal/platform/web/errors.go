package web

// ErrorResponse how we respond to clients when something goes wrong
type ErrorResponse struct {
	Error string `json:"error"`
}

// Error is used to add web information to a web request.
type Error struct {
	Err error
	Status int
}

// NewRequestError is used when a known error condition is encountered.
func NewRequestError(err error, status int) *Error {
	return &Error{
		Err: err,
		Status: status,
	}
}

// Error return string type of error message
func (e *Error) Error() string {
	return e.Err.Error()
}
