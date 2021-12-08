package handlers

import (
	"github.com/esmaeilmirzaee/grage/internal/platform/web"
	"github.com/esmaeilmirzaee/grage/internal/product"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

// ProductService is used to add database and log to a request.
type ProductService struct {
	DB *sqlx.DB
	Log *log.Logger
}

// List returns all the products stored in the database
func (p *ProductService) List(w http.ResponseWriter, r *http.Request) error {
	list, err := product.List(p.DB)
	if err != nil {
		return err
	}

	return web.Respond(w, list, http.StatusOK)
}


// Retrieve returns a product to the browser
func (p *ProductService) Retrieve(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	prod, err := product.Retrieve(p.DB, id)
	if err != nil {
		return err
	}

	return web.Respond(w, prod, http.StatusOK)
}

// Create decodes a json document from a POST request and creates a new Product.
func (p *ProductService) Create(w http.ResponseWriter, r *http.Request) error {
	var np product.NewProduct
	if err := web.Decode(r, &np); err != nil {
		return err
	}

	 prod, err := product.Create(p.DB, np, time.Now())
	 if err != nil {
		 return err
	}

	return web.Respond(w, prod, http.StatusCreated)
}
