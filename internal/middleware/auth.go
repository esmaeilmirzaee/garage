package middleware

import (
	"context"
	"github.com/esmaeilmirzaee/grage/internal/auth"
	"github.com/esmaeilmirzaee/grage/internal/platform/web"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
	"net/http"
	"strings"
)

// ErrForbidden is returned when an authenticated user does not have a
// sufficient role for an action
var ErrForbidden = web.NewRequestError(errors.New("You are not authorized for that action"), http.StatusForbidden)

// Authenticate validates a JWT from the 'Authorization' token.
func Authenticate(authenticator *auth.Authenticator) web.Middleware {
	// This is the actual middleware function to be executed.
	f := func(after web.Handler) web.Handler {
		// Wrap this handler around the next one provided.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			// Trace the application
			ctx, span := trace.StartSpan(ctx, "internal.middleware.auth")
			defer span.End()

			// Parse the authorization header. Expected header is of
			// the format <Bearer> token.
			parts := strings.Split(r.Header.Get("Authorization"), " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				err := errors.New("Expected authorization header format: <Bearer> Token")
				return web.NewRequestError(err, http.StatusUnauthorized)
			}

			// It is possible to measure just a function
			_, span = trace.StartSpan(ctx, "internal.middleware.authenticator.parseclaims")
			claims, err := authenticator.ParseClaims(parts[1])
			if err != nil {
				return web.NewRequestError(err, http.StatusUnauthorized)
			}
			span.End()

			// Add claims to the context, so they can be retrieved later.
			ctx = context.WithValue(ctx, auth.Key, claims)

			return after(ctx, w, r)
		}

		return h
	}

	return f
}

// HasRole validates that an authenticated user has at least one role from
// a specified list. This method constructs the actual function that is used.
func HasRole(roles ...string) web.Middleware {
	// This is the actual middleware function to be executed.
	f := func(after web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			claims, ok := ctx.Value(auth.Key).(auth.Claims)
			if !ok {
				return errors.New("Claims missing from context: HasRole called without/before middleware")
			}

			if !claims.HasRole(roles...) {
				return ErrForbidden
			}
			return after(ctx, w, r)
		}
		return h
	}
	return f
}
