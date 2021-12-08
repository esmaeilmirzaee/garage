package web

import (
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
)

// Respond returns the client provided data
func Respond(w http.ResponseWriter, value interface{}, statusCode int) error {
	w.Header().Set("content-type", "application/json; charset=urf8")
	w.WriteHeader(statusCode)

	data, err := json.Marshal(value)
	if err != nil {
		return errors.Wrap(err, "Could not marshal the value")
	}

	if _, err := w.Write(data); err != nil {
		return errors.Wrap(err, "Could not write to the client")
	}

	return nil
}

// RespondError knows how to handle errors going out to the client
func RespondError(w http.ResponseWriter, err error) error {
	if webErr, ok := err.(*Error); ok {
		resp := &ErrorResponse{
			Error: webErr.Err.Error(),
		}

		return Respond(w, resp, webErr.Status)
	}

	resp := ErrorResponse{
		Error: http.StatusText(http.StatusInternalServerError),
	}

	return Respond(w, resp, http.StatusInternalServerError)
}
