package middleware

import (
	"context"
	"github.com/esmaeilmirzaee/grage/internal/platform/web"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
	"net/http"
)

// Panics recovers from panics and converts the panic to an error so if it is
// reported in Metrics and handled in Errors.
func Panics() web.Middleware {
	// This is the core functionality of our middleware
	f := func(after web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {
			ctx, span := trace.StartSpan(ctx, "internal.middleware.panics")
			defer span.End()

			// Defer a function to recover from a panic and set the error (err) return
			// variable after the fact. Using the errors package will generate
			// a stack trace.
			defer func() {
				if r := recover(); r != nil {
					err = errors.Errorf("panic %v", r)
				}
			}()

			// Call the next Handler and set its return value in the error (err)
			// variable
			return after(ctx, w, r)
		}
		return h
	}
	return f
}
