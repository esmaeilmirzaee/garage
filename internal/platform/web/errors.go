package web

// ErrorResponse how we respond to clients when something goes wrong
type ErrorResponse struct {
	Error string `json:"error"`
}
