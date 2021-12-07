package product

import "time"

// Product represents the product model
type Product struct {
	ID string `db:"product_id" json:"id"`
	Name string `db:"name" json:"name"`
	Cost int `db:"cost" json:"cost"`
	Quantity int `db:"quantity" json:"quantity"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
