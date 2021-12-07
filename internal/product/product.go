package product

import (
	"github.com/jmoiron/sqlx"
	"log"
)

// List queries a database for products
func List(db *sqlx.DB) ([]Product, error) {
	list := []Product{}
	const q = "SELECT product_id, name, cost, quantity, created_at, updated_at FROM products;"

	if err := db.Select(&list, q); err != nil {
		log.Println("internal: Could not query the database", err)
		return nil, err
	}

	return list, nil
}
