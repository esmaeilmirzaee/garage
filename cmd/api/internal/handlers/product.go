package handlers

import (
	"encoding/json"
	"github.com/esmaeilmirzaee/grage/internal/product"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
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
	id := "TODO"
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
	_, err = w.Write(data)
	if err != nil {
		p.Log.Println("Could not response the result", err)
		return
	}
}
