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
	w.Write(data)
}
