package database

import (
	"github.com/jmoiron/sqlx"
	"net/url"
)

// Open creates a connection to the database
func Open() (*sqlx.DB, error) {
	q := url.Values{}
	q.Set("sslmode", "disable")
	q.Set("timezone", "utc")

	u := url.URL{
		Scheme: "postgres",
		User: url.UserPassword("pgdmn", "secret"),
		Host: "192.168.101.2:5234",
		Path: "garage",
		RawQuery: q.Encode(),
	}

	return sqlx.Open("postgres", u.String())
}