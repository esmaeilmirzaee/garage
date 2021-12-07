package api

import (
	"context"
	"github.com/esmaeilmirzaee/grage/cmd/api/internal/handlers"
	"github.com/esmaeilmirzaee/grage/internal/platform/database"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	db, err := database.Open()
	if err != nil {
		log.Fatalln("main: Could not connect to database.", err)
	}

	ps := handlers.ProductService{
		DB: db,
	}

	// Setup applications
	api := http.Server{
		Addr: "localhost:5000",
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Handler: http.HandlerFunc(ps.List),
	}

	serverErrors := make(chan error, 1)

	go func(){
		log.Println("main: Api is listening on ", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
		case err := <- serverErrors:
			log.Fatalf("main: Listening and serving: %s", err)
		case <-shutdown:
			log.Println("main: Start shutdown")
			const timeout = 5 * time.Second
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			err := api.Shutdown(ctx)
			if err != nil {
				log.Fatalf("main: Grceful shut down did not complete in %d. %s", timeout, err)
				err = api.Close()
			}

			if err != nil {
				log.Fatalf("Could not gracefully shut down server. %s", err)
			}
	}
}