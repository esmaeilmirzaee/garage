package main

import (
	"context"
	"fmt"
	"github.com/ardanlabs/conf"
	"github.com/esmaeilmirzaee/grage/cmd/api/internal/handlers"
	"github.com/esmaeilmirzaee/grage/internal/platform/database"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/ardanlabs/conf"
)

func main() {
	var cfg struct{
		DB struct{
			User	string	`conf:"default:postgres"`
			Password	string `conf:"default:secret,noprint"`
			Name	string	`conf:"default:garage"`
			Host string	`conf:"default:192.168.101.2:5234"`
			DisableTLS	bool `conf:"default:true"`
		}
	}

	// =============================================================
	// App starting
	log.Printf("main: Started.")
	defer log.Println("main: Ended.")


	// =============================================================
	// Get configuration
	if err := conf.Parse(os.Args[1:], "SALES | ", &cfg); err != nil {
		if err == conf.ErrHelpWanted{
			usage, err := conf.Usage("SALES", &cfg)
			if err != nil {
				log.Fatalf("error: generating config usage: %v", err)
			}
			fmt.Println(usage)
			return
		}
		log.Fatalf("error: Parsing config: %s.", err)
	}

	out, err := conf.String(&cfg)
	if err != nil {
		log.Fatalf("error: Generating config output. %v", err)
	}
	log.Printf("main: Config \n%v\n", out)

	// =============================================================
	// Setup dependencies
	db, err := database.Open(database.Config{
		Host: cfg.DB.Host,
		Name: cfg.DB.Name,
		User: cfg.DB.User,
		Password: cfg.DB.Password,
		DisableTLS: cfg.DB.DisableTLS,
	})
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
		Handler: http.HandlerFunc(ps.Product),
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