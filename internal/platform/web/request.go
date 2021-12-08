package web

import (
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
)

// Decode converts a JSON document received in the body of request to an interface.
func Decode(r *http.Request, val interface{}) error {

	if err := json.NewDecoder(r.Body).Decode(val); err != nil {
		return errors.Wrap(err, "Could not marshal data")
	}

	return nil
}
