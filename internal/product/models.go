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

// NewProduct is what we require from clients to make a new Product.
type NewProduct struct {
	Name	string `json:"name"`
	Cost 	int `json:"cost"`
	Quantity int `json:"quantity"`
}

// Sale represents sale model in our database.
type Sale struct {
	ID string `db:"sale_id" json:"id"`
	ProductID string `db:"product_id" json:"product_id"`
	Paid int `db:"paid" json:"paid"`
	Quantity int `db:"quantity" json:"quantity"`
	CreatedAt	time.Time `db:"created_at" json:"created_at"`
}

// NewSale is what we require from clients to make a new Sale.
type NewSale struct {
	Quantity	int `json:"quantity"`
	Paid int 	`json:"paid"`
}
