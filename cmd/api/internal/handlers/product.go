package handlers

import (
	"github.com/esmaeilmirzaee/grage/internal/platform/web"
	"github.com/esmaeilmirzaee/grage/internal/product"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
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
	list, err := product.List(r.Context(), p.DB)
	if err != nil {
		return err
	}

	return web.Respond(w, list, http.StatusOK)
}

// Retrieve returns a product to the browser
func (p *ProductService) Retrieve(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	prod, err := product.Retrieve(r.Context(), p.DB, id)
	if err != nil {
		switch err {
		case product.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case product.ErrInvalidUUID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "looking for products %q", id)
		}
	}

	return web.Respond(w, prod, http.StatusOK)
}

// Create decodes a json document from a POST request and creates a new Product.
func (p *ProductService) Create(w http.ResponseWriter, r *http.Request) error {
	var np product.NewProduct
	if err := web.Decode(r, &np); err != nil {
		return err
	}

	 prod, err := product.Create(r.Context(), p.DB, np, time.Now())
	 if err != nil {
		 return err
	}

	return web.Respond(w, prod, http.StatusCreated)
}

// ListSales returns all sales for a Product.
func (p *ProductService) ListSales(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	lists, err := product.ListSales(r.Context(), p.DB, id)
	if err != nil {
		return errors.Wrap(err, "getting sales list")
	}

	return web.Respond(w, lists, http.StatusOK)
}

// AddSale creates a new Sale for a Produce.
func (p *ProductService) AddSale(w http.ResponseWriter, r *http.Request) error {
	var ns product.NewSale
	productID := chi.URLParam(r, "id")

	if err := web.Decode(r, &ns); err != nil {
		return errors.Wrap(err, "decoding new sale")
	}

	sale, err := product.AddSale(r.Context(), p.DB, productID, ns, time.Now())
	if err != nil {
		return errors.Wrap(err, "adding new sale")
	}

	return web.Respond(w, sale, http.StatusCreated)
}
