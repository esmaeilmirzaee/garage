package main

import (
	"fmt"
	"github.com/ardanlabs/conf"
	"github.com/esmaeilmirzaee/grage/internal/platform/database"
	schema2 "github.com/esmaeilmirzaee/grage/internal/schema"
	"github.com/google/uuid"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
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
				log.Fatalln("main: generating usage %v", err)
			}
			fmt.Println(usage)
			return
		}
		log.Fatalf("error: Could not parse config %s", err)
	}

	// Setup dependencies
	db, err := database.Open(database.Config{
		Host:       cfg.DB.Host,
		Name:       cfg.DB.Name,
		User:       cfg.DB.User,
		Password:   cfg.DB.Password,
		DisableTLS: cfg.DB.DisableTLS,
	})

	if err != nil {
		log.Fatalln("main: Could not connect to database", err)
	}

	switch cfg.Args.Num(0) {
	case "migrate":
		if err := schema2.Migrate(db); err != nil {
			log.Fatalln("main: Could not migrate database.", err)
		}
		log.Println("main: Migrate is complete")
		return
	case "seed":
		if err := schema2.Seed(db); err != nil {
			log.Fatalln("main: Could not seed the database.", err)
		}
		log.Println("main: Seed is complete")
		return
	case "uuid":
		var newUUID uuid.UUID
		for i := 0; i < 10; i++ {
			newUUID = uuid.New()
			fmt.Println(newUUID)
		}
		return
	}
}
