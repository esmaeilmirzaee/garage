package schema

import "github.com/jmoiron/sqlx"

// seeds is a string constraint containing all of the queries needed to get the database
// seeded to a useful state for development.
// -------------------------------------------------------------------------------------
// Using a constant in a .go file is an easy way to ensure the queries are part of the
// compiled executable and avoids pathing issues with the working directory. It has the
// downside that it lacks syntax highlighting and may be harder to read for some cases
// compared to using .sql files. You may also consider a combined approach using a tool
// like packr or go-bindata.
// -------------------------------------------------------------------------------------
// Note that database servers besides PostgreSQL may not support running multiple queries
// as part of the same execution so this single large constant may need to be broken up.

const seeds = `
INSERT INTO products (product_id, name, cost, quantity, created_at, updated_at) VALUES ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Comic books', 50, 42, NOW(), NOW()),('d2cabc99-9c0b-4ef8-bb6a-2bb9bd380b2c', 'McDonalds toys', 75, 120, NOW(), NOW()) ON CONFLICT DO NOTHING; 

INSERT INTO sales (sale_id, product_id, paid, quantity, created_at) VALUES ('6f70b8b7-90bf-4b43-a7c7-6c3051f5c7f1',
'd2cabc99-9c0b-4ef8-bb6a-2bb9bd380b2c', 2, 100, NOW()), ('df566f1a-d511-41eb-b612-0f8a79f7cd3a', 
'd2cabc99-9c0b-4ef8-bb6a-2bb9bd380b2c', 5, 250, NOW()), ('2feb6493-ea34-49d2-b696-11928343f8d3', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 3, 225, NOW()); 

-- Create Admin and regular User with password "gophers" 
INSERT INTO users (user_id, name, email, password, roles, created_at, updated_at) VALUES('e612a422-2239-45e3-a8e0-c0c56c71454a', 'Admin Gopher', 'admin@example.com', '', '{ADMIN,USER}', NOW(), NOW()), ('6a84703c-caaf-4c94-a0a7-b131e395abdf', 'User Gopher', 'user@example.com', '', '{USER}', NOW(), NOW()) ON CONFLICT DO NOTHING;`

// Seed runs the set of seed-data queries against db. The queries are ran in a
// transaction and rolled back if any fail.
func Seed(db *sqlx.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if _, err = tx.Exec(seeds); err != nil {
		if err = tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	return tx.Commit()
}
