package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
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
		Handler: http.HandlerFunc(Echo),
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

func Echo(w http.ResponseWriter, r *http.Request) {
	random := rand.Intn(10000)
	log.Printf("echo: Starting %d", random)
	defer log.Printf("echo: Finishing %d", random)

	time.Sleep(time.Second * 3)
	fmt.Fprintf(w, "You asked for %s %s.\n", r.Method, r.URL.Path)
}
