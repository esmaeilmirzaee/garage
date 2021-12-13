package database

import (
	"context"
	"github.com/jmoiron/sqlx"
	"net/url"

	_ "github.com/lib/pq"
)

type Config struct {
	Host       string
	Name       string
	User       string
	Password   string
	DisableTLS bool
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
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.User, cfg.Password),
		Host:     cfg.Host,
		Path:     cfg.Name,
		RawQuery: q.Encode(),
	}

	return sqlx.Open("postgres", u.String())
}

// StatusCheck returns nil if it can successfully talk to
// the database. It returns a non-nil error otherwise.
func StatusCheck(ctx context.Context, db *sqlx.DB) error {
	// Run a simple query to determine connectivity. The db
	// has a "Ping" method, but it can return false-positive when it
	// was previously able to talk to the database but the database
	// has since gone away. Running this query forces a round
	// trip to the database.
	var tmp bool
	const q = `SELECT true`
	return db.QueryRowContext(ctx, q).Scan(&tmp)
}
