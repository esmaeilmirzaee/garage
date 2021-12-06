package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	/* Application starts
	 */
	log.Printf("main: Startd")
	defer log.Println("main: Finished")
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

type Product struct {
	Name string `json:"name"`
	Cost int `json:"cost"`
	Quantity int `json:"quantity"`
}

// ListProducts gets all the products
func ListProducts(w http.ResponseWriter, r *http.Request) {
	lists := []Product{
		{Name: "Comic book", Cost: 75, Quantity: 20},
		{Name: "McDonald's toy", Cost: 25, Quantity: 120},
	}

	data, err := json.MarshalIndent(lists, "", "   ");
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
