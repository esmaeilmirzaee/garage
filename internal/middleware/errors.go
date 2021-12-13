package middleware

import (
	"github.com/esmaeilmirzaee/grage/internal/platform/web"
	"log"
	"net/http"
)

// Errors handle errors coming out of the call chain. It detects normal
// application errors which are used to respond to the client in a uniform
// way. Unexpected errors (status >= 500) are logged.
func Errors(log *log.Logger) web.Middleware {
	// This is the actual middleware function to be executed.
	f := func(before web.Handler) web.Handler {

		h := func(w http.ResponseWriter, r *http.Request) error {
			// Run the handler chain and catch any propagated error.
			if err := before(w, r); err != nil {
				// log the error
				log.Printf("Error %v, ", err)

				// Respond to the error
				if err := web.RespondError(w, err); err != nil {
					return err
				}
			}
			// Return nil to indicate the error has been handler
			return nil
		}
		return h
	}

	return f
}
