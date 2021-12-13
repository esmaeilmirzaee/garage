package product

import "time"

// Product represents the product model
type Product struct {
	ID        string    `db:"product_id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Cost      int       `db:"cost" json:"cost"`
	Quantity  int       `db:"quantity" json:"quantity"`
	Sold      int       `db:"sold" json:"sold"`
	Revenue   int       `db:"revenue" json:"revenue"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// NewProduct is what we require from clients when adding a new Product.
type NewProduct struct {
	Name     string `json:"name" validate:"required"`
	Cost     int    `json:"cost" validate:"gte=0"`
	Quantity int    `json:"quantity" validate:"gte=1"`
}

// UpdateProduct defines what information may be provided to modify an
// existing Product. All fields are optional so clients can send just
// the fields they want changed. It uses pointer fields, so we can differentiate
// between a field that was not provided and a field that was provided as explicitly
// blank. Normally we do not want to use pointers to basic types, but we make
// exceptions around marshaling/unmarshaling.
type UpdateProduct struct {
	Name     *string `json:"name" validate:"omitempty"`
	Cost     *int    `json:"cost" validate:"omitempty,gte=0"`
	Quantity *int    `json:"quantity" validate:"omitempty, gte=1"`
}

// Sale represents sale model in our database.
type Sale struct {
	ID        string    `db:"sale_id" json:"id"`
	ProductID string    `db:"product_id" json:"product_id"`
	Paid      int       `db:"paid" json:"paid"`
	Quantity  int       `db:"quantity" json:"quantity"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// NewSale is what we require from clients to make a new Sale.
type NewSale struct {
	Quantity int `json:"quantity" validate:"gte=1"`
	Paid     int `json:"paid" validate:"gte=0"`
}
