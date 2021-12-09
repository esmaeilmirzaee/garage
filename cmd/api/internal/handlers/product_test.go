package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/esmaeilmirzaee/grage/internal/platform/database/databasetest"
	"github.com/esmaeilmirzaee/grage/schema"
	"github.com/google/go-cmp/cmp"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// In the Go ecosystem it's anti-pattern to have folder named tests and store
// testing files.

// TestProducts runs a series of tests to exercise Product behavior from the
// API level. The subsets all share the same database and application for
// speed and convenience. The downside is the order the tests ran matters
// and one test may break if other tests are not ran before it. If a particular
// subset needs a fresh instance of the application it can make it, or it
// should be its own Test* function.
func TestProducts(t *testing.T) {
	db, teardown := databasetest.Setup(t)
	defer teardown()

	if err := schema.Seed(db); err != nil {
		t.Fatalf("Could not seed the database. %s", err)
	}

	log := log.New(os.Stderr, "Test: ", log.LstdFlags|log.Lshortfile)

	tests := ProductTests{app: API(log, db)}

	// The following lines create subtests.
	// These tests use the same database and share the created database so their
	// changes to the database could have behavior effect on the next one.
	// The issue could be prevented by specifying or creating a new database
	// for each subtests.
	t.Run("List", tests.List)
	t.Run("ProductCRUD", tests.ProductCRUD)
}

// ProductTests holds methods for each Product subset. This type allows
// passing dependencies for test while still providing a convenient syntax
// when subsets are registered.
type ProductTests struct {
	app http.Handler
}

func (p *ProductTests) List(t *testing.T) {
	// httptest is standard package that is really helpful to test against http API
	req := httptest.NewRequest("GET", "/v1/api/products", nil)
	resp := httptest.NewRecorder()

	p.app.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("getting: expected %v, but got %v", http.StatusOK, resp.Code)
	}

	var list []map[string]interface{}

	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		t.Fatalf("decoding: %s", err)
	}

	want := []map[string]interface{} {
		{
			"id": "",
			"name": "Comic Book",
			"cost": float64(50),
			"quantity": float64(42),
			"created_at": "",
			"updated_at": "",
		},
		{
			"id":       "",
			"name":     "McDonalds Toys",
			"cost":     float64(75),
			"quantity": float64(120),
			"created_at": "",
			"updated_at": "",
		},
	}

	if diff := cmp.Diff(want, list); diff != "" {
		t.Fatalf("Response did not match expected. Diff:\n%s.", diff)
	}
}

func (p *ProductTests) ProductCRUD(t *testing.T) {
	var created map[string]interface{}

	{ // Create
		body := strings.NewReader(`{"name": "product0", "cost": 55, "quantity": 6}`)

		req := httptest.NewRequest("POST", "/v1/api/products", body)
		req.Header.Set("Content-Type", "application/json; charset=utf8;")
		resp := httptest.NewRecorder()

		p.app.ServeHTTP(resp, req)

		if http.StatusCreated != resp.Code {
			t.Fatalf("posting: expected status code %v, got %v", http.StatusCreated, resp.Code)
		}

		if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
			t.Fatalf("decoding: %s", err)
		}

		if created["id"] == "" || created["id"] == nil {
			t.Fatal("expected non-empty product id")
		}

		if created["created_at"] == "" || created["updated_at"] == nil {
			t.Fatal("expected non-empty product date created.")
		}

		want := map[string]interface{} {
			"id": created["id"],
			"name": "product0",
			"cost": float64(55),
			"quantity": float64(6),
			"created_at": created["created_at"],
			"updated_at": created["updated_at"],
		}

		if diff := cmp.Diff(want, created); diff != "" {
			t.Fatalf("Response did not match expected. Diff: %v", diff)
		}
	}

	{ // READ
		url := fmt.Sprintf("/v1/api/products/%s", created["id"])
		req := httptest.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		p.app.ServeHTTP(resp, req)

		if http.StatusOK != resp.Code {
			t.Fatalf("retrieving: expected status code %v, got %v.", http.StatusOK, resp.Code)
		}

		var fetched map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&fetched); err != nil {
			t.Fatalf("decoding: %v", err)
		}

		// Fetched product should match the one we created.
		if diff := cmp.Diff(created, fetched); diff != "" {
			t.Fatalf("Retrieved product did not match the created one. Diff: %v.", diff)
		}
	}
}
