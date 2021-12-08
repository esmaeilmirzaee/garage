package handlers

import (
	"encoding/json"
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

	data, err := json.MarshalIndent(&list, "", "   ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		p.Log.Println("Could not marshal data", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		p.Log.Println("Could not write to the browser", err)
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

	data, err := json.MarshalIndent(prod, "", "   ")
	if err != nil {
		p.Log.Println("Could not marshal the product", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	if err != nil {
		p.Log.Println("Could not response the result", err)
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

	data, err := json.MarshalIndent(prod, "", "   ")
	if err != nil {
		p.Log.Println("Could not marshal the data", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json; charset=utf8")
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(data); err != nil {
		p.Log.Println("Could not write to the client", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
