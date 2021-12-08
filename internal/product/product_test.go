package product_test

import (
	"github.com/esmaeilmirzaee/grage/internal/platform/database/databasetest"
	"github.com/esmaeilmirzaee/grage/internal/product"
	"github.com/google/go-cmp/cmp"
	"testing"
	"time"
)

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
