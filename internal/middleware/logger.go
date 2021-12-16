package middleware

import (
	"context"
	"github.com/esmaeilmirzaee/grage/internal/platform/web"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
	"log"
	"net/http"
	"time"
)

// Logger will log a line for every
func Logger(log *log.Logger) web.Middleware {
	// This is the actual middleware to be executed.
	f := func(before web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			v, ok := ctx.Value(web.KeyValues).(*web.Values)
			if !ok {
				return errors.New("Web values missing from context")
			}

			// Trace the application
			ctx, span := trace.StartSpan(ctx, "internal.middleware.logger")
			// Postpone the end to measure the entire process
			defer span.End()

			// Run the handler
			err := before(ctx, w, r)

			log.Printf("%d %s %s (%s)", v.StatusCode, r.Method, r.URL.Path, time.Since(v.Start))

			return err
		}
		return h
	}
	return f
}
