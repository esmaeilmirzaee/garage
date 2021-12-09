package product_test

import (
	"github.com/esmaeilmirzaee/grage/internal/platform/database/databasetest"
	"github.com/esmaeilmirzaee/grage/internal/product"
	"github.com/esmaeilmirzaee/grage/schema"
	"github.com/google/go-cmp/cmp"
	"testing"
	"time"
)

// TestProducts creates and restore a new Product from the testing database
func TestProducts(t *testing.T) {
	db, cleanup := databasetest.Setup(t)
	defer cleanup()


	np := product.NewProduct{
		Name: "Comic Books",
		Cost: 5,
		Quantity: 11,
	}

	now := time.Date(2021, time.December, 5, 0, 0, 0, 0, time.UTC)

	p0, err := product.Create(db, np, now)
	if err != nil {
		t.Fatalf("Could not create new product %s", err)
	}

	p1, err := product.Retrieve(db, p0.ID)
	if err != nil {
		t.Fatalf("Could not retrieve product %q %v", p0.ID, err)
	}

	if diff := cmp.Diff(p0, p1); diff != "" {
		t.Fatalf("Stored and provided product mismatch. see diff: %v\n", diff)
	}

}

// TestList uses seed to function to seed the testing database and finally checks
// the length of seed and retrieved Products
func TestList(t *testing.T) {
	db, cleanup := databasetest.Setup(t)
	defer cleanup()

	if err := schema.Seed(db); err != nil {
		t.Fatalf("Could not seed the testing database. %v", err)
	}

	ps, err := product.List(db)
	if err != nil {
		t.Fatalf("Could not retrieve data from testing database %v", err)
	}

	if exp, got := 2, len(ps); exp != got {
		t.Fatalf("Expected %v but got %v", exp, got)
	}
}
