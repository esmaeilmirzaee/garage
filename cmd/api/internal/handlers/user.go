package handlers

import (
	"context"
	"github.com/esmaeilmirzaee/grage/internal/auth"
	"github.com/esmaeilmirzaee/grage/internal/platform/web"
	"github.com/esmaeilmirzaee/grage/internal/user"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"net/http"
)

// Users holds handlers for dealing with user.
type Users struct {
	DB            *sqlx.DB
	authenticator *auth.Authenticator
}

// Token generates an authentication token for a user. The client must include
// an email and password for the request using HTTP Basic Authentication. The
// User will be identified by email and authenticated by their password.
func (u *Users) Token(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return errors.New("Web value missing from context")
	}

	email, pass, ok := r.BasicAuth()
	if !ok {
		err := errors.New("Must provide email and password for basic auth")
		return web.NewRequestError(err, http.StatusUnauthorized)
	}

	claims, err := user.Authenticate(ctx, u.DB, v.Start, email, pass)
	if err != nil {
		switch err {
		case user.ErrAuthenticationFailure:
			return web.NewRequestError(err, http.StatusUnauthorized)
		default:
			return errors.Wrap(err, "authenticating")
		}
	}
	var tkn struct {
		Token string `json:"token"`
	}
	tkn.Token, err = u.authenticator.GenerateToken(claims)
	if err != nil {
		return errors.Wrap(err, "generating token")
	}

	return web.Respond(ctx, w, tkn, http.StatusOK)
}
