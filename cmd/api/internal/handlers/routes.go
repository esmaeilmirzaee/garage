package handlers

import (
	"github.com/esmaeilmirzaee/grage/internal/middleware"
	"github.com/esmaeilmirzaee/grage/internal/platform/web"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
)

// API constructs a handler that knows about all routes
func API(log *log.Logger, db *sqlx.DB) http.Handler {
	app := web.NewApp(log, middleware.Logger(log), middleware.Errors(log), middleware.Metrics())

	c := Check{
		DB: db,
	}
	app.Handle(http.MethodGet, "/v1/api/health", c.Health)

	u := Users{
		DB: db,
	}
	app.Handle(http.MethodGet, "/v1/api/users", u.Token)

	p := ProductService{
		DB:  db,
		Log: log,
	}
	app.Handle(http.MethodGet, "/v1/api/products", p.List)
	app.Handle(http.MethodPost, "/v1/api/products", p.Create)
	app.Handle(http.MethodGet, "/v1/api/products/{id}", p.Retrieve)
	app.Handle(http.MethodPut, "/v1/api/products/{id}", p.Update)
	app.Handle(http.MethodDelete, "/v1/api/products/{id}", p.Delete)

	app.Handle(http.MethodGet, "/v1/api/products/{id}/sales", p.ListSales)

	return app
}
