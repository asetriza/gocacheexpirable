package main

import (
	"fmt"
	"time"

	"github.com/asetriza/gocacheexpirable/cache"
	"github.com/asetriza/gocacheexpirable/fetch"
	"github.com/asetriza/gocacheexpirable/fetch/httpapi"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	// db, err := sql.Open("sqlite3", "./currencies.db")
	// if err != nil {
	// 	log.Println(err)
	// }
	// defer db.Close()

	// if err := db.Ping(); err != nil {
	// 	log.Println(err)
	// }

	// fdb := database.New(db)
	// f := fetch.New(fdb)

	// values, err := f.F.FetchAll()
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Println(values)

	// value, err := f.F.Fetch(50)
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Println(value)

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
	// value, err := f1.F.Fetch(50)
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Println(value)
	// values, err := f1.F.FetchAll()
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Println(len(values))

	c, err := cache.New(f1.F)
	if err != nil {
		fmt.Println(err)
	}
	c.ReloadEvery(time.Duration(7 * time.Second))

	cache := cache.NewCacher(c)

	fmt.Println("Hello, world!", err)

	ch := make(chan struct{})
	go func(ch chan struct{}) {
		time1 := time.Now()
		value, err := cache.C.Get(2208)
		fmt.Println(time.Now().Sub(time1))
		fmt.Println("value1: ", value, "err:", err)
		ch <- struct{}{}
	}(ch)
	<-ch

	time2 := time.Now()
	value, err := cache.C.Get(2208)
	fmt.Println(time.Now().Sub(time2))
	fmt.Println("value: ", value, "err:", err)

	time3 := time.Now()
	value, err = cache.C.Get(6000)
	fmt.Println(time.Now().Sub(time3))
	fmt.Println("value: ", value, "err:", err)

	fmt.Println("Sleep zzzzzz")
	time.Sleep(time.Duration(8 * time.Second))

	time4 := time.Now()
	value, err = cache.C.Get(50)
	fmt.Println(time.Now().Sub(time4))
	fmt.Println("value: ", value, "err:", err)

	ch1 := make(chan struct{})
	go func(ch1 chan struct{}) {
		time1 := time.Now()
		value, err := cache.C.Get(2208)
		fmt.Println(time.Now().Sub(time1))
		fmt.Println("value1: ", value, "err:", err)
		ch1 <- struct{}{}
	}(ch1)
	<-ch1

	time5 := time.Now()
	value, err = cache.C.Get(50)
	fmt.Println(time.Now().Sub(time5))
	fmt.Println("last value: ", value, "err:", err)
}
