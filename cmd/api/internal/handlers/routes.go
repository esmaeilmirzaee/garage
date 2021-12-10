package handlers

import (
	"github.com/esmaeilmirzaee/grage/internal/platform/web"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
)

// API constructs a handler that knows about all routes
func API(log *log.Logger, db *sqlx.DB) http.Handler {
	p := ProductService{
		DB: db,
		Log: log,
	}

	app := web.NewApp(log)

	app.Handle(http.MethodGet, "/v1/api/products", p.List)
	app.Handle(http.MethodPost, "/v1/api/products", p.Create)
	app.Handle(http.MethodGet, "/v1/api/products/{id}", p.Retrieve)

	app.Handle(http.MethodGet, "/v1/api/products/{id}/sales", p.ListSales)
	app.Handle(http.MethodPost, "/v1/api/products/{id}/sales", p.AddSale)

	return app
}
