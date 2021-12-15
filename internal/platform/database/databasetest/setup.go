package databasetest

import (
	"github.com/esmaeilmirzaee/grage/internal/platform/database"
	"github.com/esmaeilmirzaee/grage/internal/schema"
	"github.com/jmoiron/sqlx"
	"log"
	"testing"
	"time"
)

// Setup creates a test database inside a Docker container. It creates the
// required table structure but the database is otherwise empty.
//
// It does not return error at this intended for testing only. Instead, it will
// call Fatal on the provided testing, if anything goes wrong.
//
// It returns the database to use as well as a function to call at the end of
// the test.
func Setup(t *testing.T) (*sqlx.DB, func()) {
	t.Helper()

	c := startContainer(t)

	db, err := database.Open(database.Config{
		User:       "postgres",
		Password:   "postgres",
		Name:       "garage-testing",
		Host:       c.Host,
		DisableTLS: true,
	})
	if err != nil {
		log.Fatalf("Opening database connection %s.", err)
	}

	t.Log("Waiting for database to be ready")

	// Wait for the database to be ready. Wait 100ms longer between each attempt.
	// Do not try more than 20 times.
	var pingError error
	maxAttempts := 20
	for attempts := 1; attempts <= maxAttempts; attempts++ {
		pingError = db.Ping()
		if pingError != nil {
			break
		}
		time.Sleep(time.Duration(attempts) * 100 * time.Millisecond)
	}

	if pingError != nil {
		stopContainer(t, c)
		t.Fatalf("Waiting for the database to be ready: %v", pingError)
	}

	if err := schema.Migrate(db); err != nil {
		stopContainer(t, c)
		t.Fatalf("Migrating: %s", err)
	}

	// teardown is a function that should be invoked when the caller is done
	// with the database.
	teardown := func() {
		t.Helper()
		if err := db.Close(); err != nil {
			log.Printf("Could not close the database %v", err)
		}
		stopContainer(t, c)
	}

	return db, teardown
}
