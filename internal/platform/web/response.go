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

	if statusCode == http.StatusNoContent {
		return nil
	}

	data, err := json.Marshal(value)
	if err != nil {
		return errors.Wrap(err, "Could not marshal the value")
	}

	if _, err := w.Write(data); err != nil {
		return errors.Wrap(err, "Could not write to the client")
	}

	return nil
}

// RespondError sends an error response back to the client.
func RespondError(w http.ResponseWriter, err error) error {
	// If the error was of the type *Error, the handler has
	// a specific status code and error to return.
	if webErr, ok := errors.Cause(err).(*Error); ok {
		er := ErrorResponse{
			Error:  webErr.Err.Error(),
			Fields: webErr.Fields,
		}
		if err := Respond(w, er, webErr.Status); err != nil {
			return err
		}
		return nil
	}

	// If not, the handler sent any aribitrary error value so use 500
	er := ErrorResponse{
		Error: http.StatusText(http.StatusInternalServerError),
	}
	if err := Respond(w, er, http.StatusInternalServerError); err != nil {
		return err
	}
	return nil
}
