package product

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"log"
	"time"
)

// Predefined errors for known failure scenario.
var (
	ErrNotFound    = errors.New("Not found")
	ErrInvalidUUID = errors.New("Invalid ID")
)

// List queries a database for products
func List(ctx context.Context, db *sqlx.DB) ([]Product, error) {
	var list []Product
	const q = `SELECT p.product_id, p.name, p.cost, p.quantity, COALESCE(SUM(s.quantity), 0) AS sold, 
COALESCE(SUM(s.paid), 0) AS revenue, p.created_at, 
p.updated_at FROM products AS p LEFT JOIN sales AS s on p.product_id = s.product_id GROUP BY p.product_id;`

	if err := db.SelectContext(ctx, &list, q); err != nil {
		log.Println("internal: Could not query the database", err)
		return nil, err
	}

	return list, nil
}

// Retrieve returns a product
func Retrieve(ctx context.Context, db *sqlx.DB, id string) (*Product, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrInvalidUUID
	}

	var p Product
	q := `SELECT p.product_id, p.name, p.cost, p.quantity, COALESCE(SUM(s.paid), 0) AS revenue, 
COALESCE(SUM(s.quantity), 0) AS sold, 
p.created_at, 
p.updated_at FROM products AS p LEFT JOIN sales AS s ON s.product_id = p.product_id WHERE p.product_id = $1 GROUP BY p.
product_id`

	if err := db.GetContext(ctx, &p, q, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &p, nil
}

// Create makes a new Product.
func Create(ctx context.Context, db *sqlx.DB, np NewProduct, now time.Time) (*Product, error) {
	p := Product{
		ID:        uuid.New().String(),
		Name:      np.Name,
		Cost:      np.Cost,
		Quantity:  np.Quantity,
		CreatedAt: now,
		UpdatedAt: now,
	}

	q := `INSERT INTO products (product_id, name, cost, quantity, created_at, updated_at) VALUES($1, $2, $3, $4, 
$5, $6)`

	if _, err := db.ExecContext(ctx, q, p.ID, p.Name, p.Cost, p.Quantity, p.CreatedAt, p.UpdatedAt); err != nil {
		return nil, errors.Wrap(err, "Cannot create a new product")
	}

	return &p, nil
}

// Update modifies data about a Product. It will error if the specified
// ID is invalid or does not reference an existing Product.
func Update(ctx context.Context, db *sqlx.DB, id string, update UpdateProduct, now time.Time) error {

	p, err := Retrieve(ctx, db, id)
	if err != nil {
		return err
	}

	if update.Name != nil {
		p.Name = *update.Name
	}

	if update.Cost != nil {
		p.Cost = *update.Cost
	}

	if update.Quantity != nil {
		p.Quantity = *update.Quantity
	}

	p.UpdatedAt = now

	const q = `UPDATE products SET "name" = $2, "cost" = $3, "quantity" = $4, "updated_at" = $5 WHERE product_id = $1;`

	_, err = db.ExecContext(ctx, q, id, p.Name, p.Cost, p.Quantity, p.UpdatedAt)
	if err != nil {
		return errors.Wrap(err, "Updating failed")
	}

	return nil
}

// Delete removes a Product
func Delete(ctx context.Context, db *sqlx.DB, ProductID string) error {
	if _, err := uuid.Parse(ProductID); err != nil {
		return ErrInvalidUUID
	}
	const q = `DELETE FROM products WHERE product_id = $1`

	if _, err := db.ExecContext(ctx, q, ProductID); err != nil {
		return errors.Wrapf(err, "deleting a product %q", ProductID)
	}

	return nil
}

// AddSale creates a new Sale.
func AddSale(ctx context.Context, db *sqlx.DB, ProductID string, ns NewSale,
	now time.Time) (*NewSale,
	error) {
	if _, err := uuid.Parse(ProductID); err != nil {
		return nil, ErrInvalidUUID
	}

	s := Sale{
		ID:        uuid.New().String(),
		ProductID: ProductID,
		Paid:      ns.Paid,
		Quantity:  ns.Quantity,
		CreatedAt: now,
	}

	const q = `INSERT INTO sales (sale_id, product_id, paid, quantity, created_at) VALUES ($1, $2, $3, $4, $5);`

	if _, err := db.ExecContext(ctx, q, s.ID, s.ProductID, s.Paid, s.Quantity, s.CreatedAt); err != nil {
		return nil, errors.Wrap(err, "Could not create new sale")
	}

	return &ns, nil
}

// ListSales returns all sales for a Product.
func ListSales(ctx context.Context, db *sqlx.DB, ProductID string) ([]Sale, error) {
	if _, err := uuid.Parse(ProductID); err != nil {
		return nil, ErrInvalidUUID
	}

	q := `SELECT product_id, sale_id, paid, quantity, created_at FROM sales WHERE product_id = $1;`
	var list []Sale

	if err := db.SelectContext(ctx, &list, q, ProductID); err != nil {
		return nil, errors.Wrap(err, "Could not query the database")
	}

	return list, nil
}
