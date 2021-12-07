package main

import (
	"flag"
	"github.com/esmaeilmirzaee/grage/internal/platform/database"
	"github.com/esmaeilmirzaee/grage/schema"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	// Setup dependencies
	db, err := database.Open()
	if err != nil {
		log.Fatalln("main: Could not connect to database", err)
	}

	flag.Parse()
	switch flag.Arg(0) {
	case "migrate":
		if err := schema.Migrate(db); err != nil {
			log.Fatalln("main: Could not migrate database.", err)
		}
		log.Println("main: Migrate is complete")
		return
	case "seed":
		if err := schema.Seed(db); err != nil {
			log.Fatalln("main: Could not seed the database.", err)
		}
		log.Println("main: Seed is complete")
		return
	}
}
