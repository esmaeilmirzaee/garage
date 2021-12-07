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

const seeds = `TRUNCATE TABLE products; INSERT INTO products (product_id, name, cost, quantity, createdAt, 
updatedAt) VALUES ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Comic books', 50, 42, '2021-01-12 00:00:01.000001+00', '2021-01-12 00:00:01.000001+00'),('d2cabc99-9c0b-4ef8-bb6a-2bb9bd380b2c', 'McDonalds toys', 75, 120, '2021-01-12 00:00:01.000001+00', '2021-01-12 00:00:01.000001+00') ON CONFLICT DO NOTHING;`

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
















