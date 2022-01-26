package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/asetriza/gocacheexpirable/cache"
	"github.com/asetriza/gocacheexpirable/fetch"
	"github.com/asetriza/gocacheexpirable/fetch/database"
	"github.com/asetriza/gocacheexpirable/fetch/httpapi"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	c := cache.New()
	fmt.Println("Hello, world!", c)

	db, err := sql.Open("sqlite3", "./currencies.db")
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Println(err)
	}

	fdb := database.New(db)
	f := fetch.New(fdb)

	values, err := f.F.FetchAll()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(values)

	value, err := f.F.Fetch(50)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(value)

	api := httpapi.HttpAPI{
		Scheme: "https",
		Host:   "api.skypicker.com",
		Endpoints: httpapi.Endpoints{
			Location: httpapi.Endpoint{
				Path:   "/locations",
				Method: "GET",
			},
			Locations: httpapi.Endpoint{
				Path:   "/locations/graphql",
				Method: "GET",
			},
		},
		Timeout: time.Duration(5 * time.Second),
	}

	f1 := fetch.New(api)
	value, err = f1.F.Fetch(2208)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(value)
	values, err = f1.F.FetchAll()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(len(values))
}
