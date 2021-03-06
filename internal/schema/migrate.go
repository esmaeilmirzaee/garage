package schema

import (
	"github.com/GuiaBolso/darwin"
	"github.com/jmoiron/sqlx"
)

// Migrations contains the queries needed to construct the database schema.
// Entries should never be removed from this slice once they have been ran in production.
//
// Including the queries directly in this file has the same pros/cons mentioned in seed.go
var migrations = []darwin.Migration{
	{
		Version:     1,
		Description: "Add products",
		Script:      `CREATE TABLE products (product_id UUID, name TEXT, cost INT, quantity INT, created_at TIMESTAMP, updated_at TIMESTAMP, PRIMARY KEY (product_id));`,
	},
	{
		Version:     2,
		Description: "Create sales tables",
		Script: `CREATE TABLE sales (sale_id UUID, product_id UUID, paid INT, quantity INT, 
created_at TIMESTAMP,   PRIMARY KEY (sale_id), FOREIGN KEY (product_id) REFERENCES products(product_id) ON DELETE CASCADE);`,
	},
	{
		Version:     3,
		Description: "Create users table",
		Script: `CREATE TABLE IF NOT EXISTS users (user_id UUID, name TEXT, email TEXT UNIQUE, password TEXT, 
roles TEXT[], created_at TIMESTAMP, updated_at TIMESTAMP, PRIMARY KEY (user_id));`,
	},
	{
		Version:     4,
		Description: "Add user column to products table",
		Script:      `ALTER TABLE products ADD COLUMN user_id UUID DEFAULT '00000000-0000-0000-0000-000000000000';`,
	},
}

// Migrate attempts to bring the schema for db up to date with the migrations
// defined in this package.
func Migrate(db *sqlx.DB) error {
	driver := darwin.NewGenericDriver(db.DB, darwin.PostgresDialect{})

	d := darwin.New(driver, migrations, nil)

	return d.Migrate()
}
