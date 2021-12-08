package product

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"log"
	"time"
)

// Predefined errors for known failure scenario.
var (
	ErrNotFound = errors.New("Not found")
	ErrInvalidUUID = errors.New("Invalid ID")
)

// List queries a database for products
func List(db *sqlx.DB) ([]Product, error) {
	var list []Product
	const q = "SELECT product_id, name, cost, quantity, created_at, updated_at FROM products;"

	if err := db.Select(&list, q); err != nil {
		log.Println("internal: Could not query the database", err)
		return nil, err
	}

	return list, nil
}

// Retrieve returns a product
func Retrieve(db *sqlx.DB, id string) (*Product, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrInvalidUUID
	}

	var p Product
	q := `SELECT product_id, name, cost, quantity, created_at, updated_at FROM products WHERE product_id = $1`

	if err := db.Get(&p, q, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &p, nil
}

// Create makes a new Product.
func Create(db *sqlx.DB, np NewProduct, now time.Time) (*Product, error) {
	log.Println(uuid.New(), uuid.NewString())
	p := Product{
		ID: uuid.New().String(),
		Name: np.Name,
		Cost: np.Cost,
		Quantity: np.Quantity,
		CreatedAt: now,
		UpdatedAt: now,
	}

	q := `INSERT INTO products (product_id, name, cost, quantity, created_at, updated_at) VALUES($1, $2, $3, $4, 
$5, $6)`

	if _, err := db.Exec(q, p.ID, p.Name, p.Cost, p.Quantity, p.CreatedAt, p.UpdatedAt); err != nil {
		return nil, errors.Wrap(err, "Cannot create a new product")
	}

	return &p, nil
}
