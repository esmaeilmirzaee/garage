package handlers

import (
	"encoding/json"
	"github.com/esmaeilmirzaee/grage/internal/product"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
)

type ProductService struct {
	DB *sqlx.DB
}

func (p *ProductService) Product(w http.ResponseWriter, r *http.Request) {
	list, err := product.List(p.DB)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("handlers: Could not receive database data.")
		return
	}

	data, err := json.MarshalIndent(&list, "", "   ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Could not marshal data", err)
		return
	}

	w.Write(data)
}