package database

import (
	"github.com/jmoiron/sqlx"
	"net/url"

	_ "github.com/lib/pq"
)

type Config struct {
	Host		string
	Name		string
	User		string
	Password	string
	DisableTLS	bool
}

// Open knows how to open a database connection
func Open(cfg Config) (*sqlx.DB, error) {
	q := url.Values{}

	q.Set("timezone", "utc")
	// sslmode only "require" (default), "verify-full", "verify-ca", and "disable" supported
	q.Set("sslmode", "require")
	if cfg.DisableTLS {
		q.Set("sslmode", "disable")
	}

	u := url.URL{
		Scheme: "postgres",
		User: url.UserPassword(cfg.User, cfg.Password),
		Host: cfg.Host,
		Path: cfg.Name,
		RawQuery: q.Encode(),
	}

	return sqlx.Open("postgres", u.String())
}