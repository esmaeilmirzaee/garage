package handlers

import (
	"encoding/json"
	"github.com/esmaeilmirzaee/grage/internal/platform/web"
	"github.com/esmaeilmirzaee/grage/internal/product"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

// ProductService os
type ProductService struct {
	DB *sqlx.DB
	Log *log.Logger
}

// List returns all the products stored in the database
func (p *ProductService) List(w http.ResponseWriter, r *http.Request) {
	list, err := product.List(p.DB)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		p.Log.Println("handlers: Could not receive database data.", err)
		return
	}

	if err = web.Response(w, list, http.StatusOK); err != nil {
		p.Log.Println("Could not response to the client", err)
		return
	}
}


// Retrieve returns a product to the browser
func (p *ProductService) Retrieve(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	prod, err := product.Retrieve(p.DB, id)
	if err != nil {
		p.Log.Println("Could not receive the product", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = web.Response(w, prod, http.StatusOK); err != nil {
		p.Log.Println("Could not response to the client", err)
		return
	}
}

// Create decodes a json document from a POST request and creates a new Product.
func (p *ProductService) Create(w http.ResponseWriter, r *http.Request) {
	var np product.NewProduct
	if err := json.NewDecoder(r.Body).Decode(&np); err != nil {
		p.Log.Println("Could not decode product", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	 prod, err := product.Create(p.DB, np, time.Now())
	 if err != nil {
		p.Log.Println("Could not store in the database", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := web.Response(w, prod, http.StatusCreated); err != nil {
		p.Log.Println("Response to the user failed", err)
		return
	}
}
