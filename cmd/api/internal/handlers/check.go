package handlers

import (
	"github.com/esmaeilmirzaee/grage/internal/platform/database"
	"github.com/esmaeilmirzaee/grage/internal/platform/web"
	"github.com/jmoiron/sqlx"
	"net/http"
)

// Check has handlers to implement service orchestration
type Check struct {
	DB *sqlx.DB

	// ADD OTHER STATE LIKE THE LOGGER IF NEEDED
}

// Health responds with a 200 OK if the service is healthy
// and ready for the traffic.
func (c *Check) Health(w http.ResponseWriter, r *http.Request) error {
	var health struct {
		Status string `json:"status"`
	}

	// Check if the database is ready
	if err := database.StatusCheck(r.Context(), c.DB); err != nil {
		// If the database is not ready we will tell the client
		// and use a 500 status. Do not respond by just returning
		// an error because further up in the call stack will
		// interpret that as an unhandled error.
		health.Status = "db not ready"
		return web.Respond(r.Context(), w, health, http.StatusInternalServerError)
	}

	health.Status = "OK"
	return web.Respond(r.Context(), w, health, http.StatusOK)
}
