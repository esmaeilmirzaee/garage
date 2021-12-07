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
	q.Set("sslmode", "required")
	if cfg.DisableTLS {
		q.Set("sslmode", "disabled")
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