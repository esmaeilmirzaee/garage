package middleware

import (
	"github.com/esmaeilmirzaee/grage/internal/platform/web"
	"log"
	"net/http"
	"time"
)

// Logger will log a line for every
func Logger(log *log.Logger) web.Middleware {
	// This is the actual middleware to be executed.
	f := func(before web.Handler) web.Handler {
		h := func(w http.ResponseWriter, r *http.Request) error {
			start := time.Now()
			// Run the handler
			err := before(w, r)

			log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))

			return err
		}
		return h
	}
	return f
}
