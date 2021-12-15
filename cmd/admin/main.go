package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/ardanlabs/conf"
	"github.com/esmaeilmirzaee/grage/internal/auth"
	"github.com/esmaeilmirzaee/grage/internal/platform/database"
	"github.com/esmaeilmirzaee/grage/internal/schema"
	"github.com/esmaeilmirzaee/grage/internal/user"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	if err := run(); err != nil {
		log.Printf("error %s", err)
		os.Exit(1)
	}
}

func run() error {
	// =============================================================
	// Configuration
	var cfg struct {
		DB struct {
			Host       string `conf:"default:192.168.101.2:5234"`
			Name       string `conf:"default:garage"`
			User       string `conf:"default:pgdmn"`
			Password   string `conf:"default:secret"`
			DisableTLS bool   `conf:"default:true"`
		}
		Args conf.Args
	}

	if err := conf.Parse(os.Args[1:], "SALES | ", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage("SALES | ", &cfg)
			if err != nil {
				return errors.Wrap(err, "main: generating usage")
			}
			fmt.Println(usage)
			return nil
		}
		return errors.Wrap(err, "error: parsing")
	}

	// This is used for multiple commands below.
	dbConfig := database.Config{
		Host:       cfg.DB.Host,
		Name:       cfg.DB.Name,
		User:       cfg.DB.User,
		Password:   cfg.DB.Password,
		DisableTLS: cfg.DB.DisableTLS,
	}

	var err error
	switch cfg.Args.Num(0) {
	case "migrate":
		err = migrate(dbConfig)
	case "seed":
		err = seed(dbConfig)
	case "useradd":
		err = useradd(dbConfig, cfg.Args.Num(1))
	case "keygen":
		err = keygen(cfg.Args.Num(1))
	case "uuid":
		var newUUID uuid.UUID
		for i := 0; i < 10; i++ {
			newUUID = uuid.New()
			fmt.Println(newUUID)
		}
		return nil
	default:
		err = errors.New("Must specify a command")
	}

	if err != nil {
		return err
	}

	return nil
}

func migrate(dbConfig database.Config) error {
	db, err := database.Open(dbConfig)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := schema.Migrate(db); err != nil {
		return err
	}

	fmt.Println("Migrating is complete")
	return nil
}

func seed(dbConfig database.Config) error {
	db, err := database.Open(dbConfig)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := schema.Seed(db); err != nil {
		return err
	}

	fmt.Println("Seeding is complete")
	return nil
}

func useradd(dbConfig database.Config, email string) error {
	db, err := database.Open(dbConfig)
	if err != nil {
		return err
	}
	defer db.Close()

	if email == "" {
		return errors.New("User creation must be called with an additional argument | email")
	}

	fmt.Print("Please enter password: ")
	var password string
	if _, err := fmt.Scanf("%v\n", &password); err != nil {
		return errors.Wrap(err, "entering password")
	}
	if password == "" {
		fmt.Println("Canceling")
		return nil
	}

	fmt.Printf("Admin user will be created with email %q", email)
	fmt.Printf("Continue? (Y/N)")
	var confirm byte
	if _, err := fmt.Scanf("%c\n", &confirm); err != nil {
		return errors.Wrap(err, "processing response")
	}

	if string(confirm) != "y" || string(confirm) != "Y" && string(confirm) == "n" || string(confirm) == "N" {
		fmt.Println("Canceling")
		return nil
	}

	ctx := context.Background()
	nu := user.NewUser{
		Name:     email,
		Password: password,
		Email:    email,
		Roles:    []string{auth.RoleAdmin, auth.RoleUser},
	}

	user, err := user.Create(ctx, db, nu, time.Now())
	if err != nil {
		return err
	}

	fmt.Printf("User created %q", user.ID)

	return nil
}

// keygen creates an x509 private key for signing auth token.
func keygen(path string) error {
	if path == "" {
		return errors.New("keygen missing argument for key path")
	}

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return errors.Wrap(err, "generating keys")
	}

	file, err := os.Create(path)
	if err != nil {
		return errors.Wrap(err, "Creating private file")
	}
	defer file.Close()

	block := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	if err := pem.Encode(file, &block); err != nil {
		return errors.Wrap(err, "Encoding to private key")
	}

	if err := file.Close(); err != nil {
		return errors.Wrap(err, "Closing private file")
	}

	return nil
}
