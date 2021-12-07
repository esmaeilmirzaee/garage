package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/esmaeilmirzaee/grage/schema"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	/* Application starts
	 */
	log.Printf("main: Startd")
	defer log.Println("main: Finished")
	// =========================================================
	// Setup dependencies
	log.Println("Setup database connection.")
	db, err := openDB()
	if err != nil{
		log.Println("Could not connect to database")
		log.Fatal(err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Println(err)
		}
	}()

	// Handling migration and seed requests
	flag.Parse()
	switch flag.Arg(0) {
	case "migrate":
		if err := schema.Migrate(db); err != nil {
			log.Fatal("Failed to migrate database", err)
		}
		log.Println("Migrate is complete.")
		return
	case "seed":
		if err := schema.Seed(db); err != nil {
			log.Fatal("Failed to seed the database", err)
		}
		log.Println("Seed is complete")
		return
	}



	// =========================================================
	api := http.Server{
		Addr: "localhost:5000",
		ReadTimeout: time.Second * 5,
		WriteTimeout: time.Second * 5,
		Handler: http.HandlerFunc(ListProducts),
	}

	// Make a channel to listen for errors coming from listener, use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serviceErrors := make(chan error, 1)

	// Start the service and listen to the requests
	go func(){
		log.Printf("main: Api is listening on %s", api.Addr)
		serviceErrors <- api.ListenAndServe()
	}()

	// Make a channel to listen for an interrupt or terminate signal from the os.
	// Use a buffered channel because the signal package requires.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// =========================================================
	// Shutdown
	// Blocking main and waiting for shutdown
	select {
	case err := <- serviceErrors:
		log.Fatalf("main: Listening and Serving: %s", err)

	case <-shutdown:
		log.Printf("main: Start shutdown")

		// Give outstanding requests a deadline for completion.
		const timeout = 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		// Asking listener to shut down and load shed.
		err := api.Shutdown(ctx)
		if err != nil {
			log.Printf("main: Graceul shutdown did not complete in %v : %v.", timeout, err)
			err = api.Close()
		}

		if err != nil {
			log.Fatalf("main: Could not shotdown server gracefully %v", err)
		}
	}
}

func openDB() (*sqlx.DB, error) {
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

type Product struct {
	Name string `json:"name"`
	Cost int `json:"cost"`
	Quantity int `json:"quantity"`
}

// ListProducts gets all the products
func ListProducts(w http.ResponseWriter, r *http.Request) {
	lists := []Product{}

	lists = append(lists, Product{Name: "Comic book", Cost: 75, Quantity: 20})
	lists =append(lists, Product{Name: "McDonald's toy", Cost: 25, Quantity: 120})


	

	data, err := json.MarshalIndent(lists, "", "   ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Cannot generate json object")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(data); err != nil {
		log.Println("Cannot respond to the user")
	}
}
