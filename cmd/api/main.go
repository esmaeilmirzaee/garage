package main

import (
	"context"
	"contrib.go.opencensus.io/exporter/zipkin"
	"crypto/rsa"
	"fmt"
	"github.com/ardanlabs/conf"
	"github.com/dgrijalva/jwt-go"
	"github.com/esmaeilmirzaee/grage/cmd/api/internal/handlers"
	"github.com/esmaeilmirzaee/grage/internal/auth"
	"github.com/esmaeilmirzaee/grage/internal/platform/database"
	"go.opencensus.io/trace"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"

	_ "github.com/ardanlabs/conf"

	_ "expvar"         // Register the /debug/vars handler | metric middleware
	_ "net/http/pprof" // Register the /debug/pprof handler | Profiling middleware

	// Register OpenZipkin to send span
	_ "contrib.go.opencensus.io/exporter/zipkin"
	openzipkin "github.com/openzipkin/zipkin-go"
	zipkinHTTP "github.com/openzipkin/zipkin-go/reporter/http"
	_ "go.opencensus.io/trace"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	log := log.New(os.Stdout, "SALES | ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	var cfg struct {
		Web struct {
			Address         string        `conf:"default:localhost:5000"`
			Debug           string        `conf:"default:localhost:6060"`
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:5s"`
			ShutdownTimeout time.Duration `conf:"default:5s"`
		}
		DB struct {
			User       string `conf:"default:pgdmn"`
			Password   string `conf:"default:secret,noprint"`
			Name       string `conf:"default:garage"`
			Host       string `conf:"default:192.168.101.2:5234"`
			DisableTLS bool   `conf:"default:true"`
		}
		Auth struct {
			PrivateKeyFile string `conf:"default:1"`
			KeyID          string `conf:"default:private.pem"`
			Algorithm      string `conf:"default:RS256"`
		}
		Trace struct {
			URL         string  `conf:"default:http://192.168.101.2:9411/api/v2/spans"`
			Service     string  `conf:"default:grage-api"`
			Probability float64 `conf:"default:1"`
		}
	}

	// =============================================================
	// App starting
	log.Printf("main: Started.")
	defer log.Println("main: Ended.")

	// =============================================================
	// Get configuration
	if err := conf.Parse(os.Args[1:], "SALES | ", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage("SALES", &cfg)
			if err != nil {
				return errors.Wrap(err, "generating config usage")
			}
			fmt.Println(usage)
			return nil
		}
		return errors.Wrap(err, "Parsing config.")
	}

	out, err := conf.String(&cfg)
	if err != nil {
		return errors.Wrap(err, "Generating config output.")
	}
	log.Printf("main: Config \n%v\n", out)

	// =============================================================
	// Initialize authentication support
	authenticator, err := createAuth(cfg.Auth.PrivateKeyFile, cfg.Auth.KeyID, cfg.Auth.Algorithm)
	if err != nil {
		return errors.Wrap(err, "constructing authenticator")
	}

	// =============================================================
	// Setup dependencies
	// Start database
	db, err := database.Open(database.Config{
		Host:       cfg.DB.Host,
		Name:       cfg.DB.Name,
		User:       cfg.DB.User,
		Password:   cfg.DB.Password,
		DisableTLS: cfg.DB.DisableTLS,
	})
	if err != nil {
		return errors.Wrap(err, "Could not connect to database.")
	}

	// =============================================================
	// Start tracing session
	closer, err := registerTracer(cfg.Trace.Service, cfg.Web.Address, cfg.Trace.URL, cfg.Trace.Probability)
	if err != nil {
		return err
	}
	defer closer()

	// =============================================================
	// Start Debug Service
	go func() {
		log.Println("main: Profile API is listening on %v", cfg.Web.Debug)
		err := http.ListenAndServe(cfg.Web.Debug, http.DefaultServeMux)
		log.Println("main: Debug service ended %v", err)
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Start API Service
	api := http.Server{
		Addr:         cfg.Web.Address,
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		Handler:      handlers.API(shutdown, log, db, authenticator),
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Println("main: Api is listening on ", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		return errors.Wrap(err, "Listening and serving")
	case sig := <-shutdown:
		log.Println("main: Start shutdown", sig)
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		err := api.Shutdown(ctx)
		if err != nil {
			return errors.Wrap(err, "main: Grceful shut down did not complete.")
		}

		if err != nil {
			return errors.Wrap(err, "Could not gracefully shut down server.")
		}

		if sig == syscall.SIGSTOP {
			return errors.New("Integrity error detected. Asking for sel shutdown")
		}
	}

	return nil
}

// createAuth creates an x509 private key.
func createAuth(privateKeyFile, keyID, algorithm string) (*auth.Authenticator, error) {
	keyContents, err := ioutil.ReadFile(privateKeyFile)
	if err != nil {
		return nil, errors.Wrap(err, "reading auth private key")
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(keyContents)
	if err != nil {
		return nil, errors.Wrap(err, "passing auth private key")
	}

	public := auth.NewSimpleKeyLookup(keyID, key.Public().(*rsa.PublicKey))
	return auth.NewAuthenticator(key, keyID, algorithm, public)
}

// registerTracer registers for a zipkin tracer
// probability is a percentage of requests that should be monitored
// 1 equals 100%; or all the requests and 0.1 means 10%
func registerTracer(service, httpAddr, traceURL string, probability float64) (func() error, error) {
	localEndPoint, err := openzipkin.NewEndpoint(service, httpAddr)
	if err != nil {
		return nil, errors.Wrap(err, "Creating local endpoint zipkin")
	}
	reporter := zipkinHTTP.NewReporter(traceURL)

	trace.RegisterExporter(zipkin.NewExporter(reporter, localEndPoint))
	trace.ApplyConfig(trace.Config{
		DefaultSampler: trace.ProbabilitySampler(probability),
	})

	return reporter.Close, nil
}
