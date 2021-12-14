package middleware

import (
	"github.com/esmaeilmirzaee/grage/internal/platform/web"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"time"
)

// Logger will log a line for every
func Logger(log *log.Logger) web.Middleware {
	// This is the actual middleware to be executed.
	f := func(before web.Handler) web.Handler {
		h := func(w http.ResponseWriter, r *http.Request) error {
			v, ok := r.Context().Value(web.KeyValues).(*web.Values)
			if !ok {
				return errors.New("Web values missing from context")
			}

			// Run the handler
			err := before(w, r)

			log.Printf("%d %s %s (%s)", v.StatusCode, r.Method, r.URL.Path, time.Since(v.Start))

			return err
		}
		return h
	}
	return f
}
