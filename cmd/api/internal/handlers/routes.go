package handlers

import (
	"github.com/esmaeilmirzaee/grage/internal/auth"
	"github.com/esmaeilmirzaee/grage/internal/middleware"
	"github.com/esmaeilmirzaee/grage/internal/platform/web"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
)

// API constructs a handler that knows about all routes
func API(log *log.Logger, db *sqlx.DB, authenticator *auth.Authenticator) http.Handler {
	// It is almost impossible to put auth middleware here because it would block
	// all the routes; even the authentication mechanism
	app := web.NewApp(log, middleware.Logger(log), middleware.Errors(log), middleware.Metrics())

	c := Check{
		DB: db,
	}
	app.Handle(http.MethodGet, "/v1/api/health", c.Health)

	u := Users{
		DB:            db,
		authenticator: authenticator,
	}
	app.Handle(http.MethodGet, "/v1/api/users", u.Token)

	p := ProductService{
		DB:  db,
		Log: log,
	}
	// the following routes require authorizations
	app.Handle(http.MethodGet, "/v1/api/products", p.List, middleware.Authenticate(authenticator))
	app.Handle(http.MethodPost, "/v1/api/products", p.Create, middleware.Authenticate(authenticator))
	app.Handle(http.MethodGet, "/v1/api/products/{id}", p.Retrieve, middleware.Authenticate(authenticator))
	app.Handle(http.MethodPut, "/v1/api/products/{id}", p.Update, middleware.Authenticate(authenticator))
	app.Handle(http.MethodDelete, "/v1/api/products/{id}", p.Delete, middleware.Authenticate(authenticator))

	app.Handle(http.MethodGet, "/v1/api/products/{id}/sales", p.ListSales, middleware.Authenticate(authenticator))
	app.Handle(http.MethodPost, "/v1/api/products/{id}/sales", p.AddSale, middleware.Authenticate(authenticator))

	return app
}
