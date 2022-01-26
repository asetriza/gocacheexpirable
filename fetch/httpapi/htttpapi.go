package httpapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const defaultLtype string = "general"
const defaultKey string = "int_id"
const defaultQuery string = "{dump (options: {location_types: [\"airport\"], active_only: \"true\"}) {int_id id}}"

type HttpAPI struct {
	Scheme    string
	Host      string
	Endpoints Endpoints
	Timeout   time.Duration
}

type Endpoint struct {
	Path   string
	Method string
}

type Endpoints struct {
	Locations Endpoint
	Location  Endpoint
}

func (ha HttpAPI) Fetch(id int) (string, error) {

	method, url := ha.locationURL(defaultLtype, defaultKey, fmt.Sprint(id))
	request, err := http.NewRequest(method, fmt.Sprint(url), bytes.NewBuffer([]byte{}))
	if err != nil {
		log.Println(err)
		return "", err
	}
	request.Header.Set("Content-type", "application/json")

	client := http.Client{
		Timeout: ha.Timeout,
	}
	resp, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer resp.Body.Close()

	var location LocationBody
	err = json.Unmarshal(body, &location)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return fmt.Sprint(location.Locations[0].IntID), nil
}

func (ha HttpAPI) FetchAll() (map[int]string, error) {

	method, url := ha.locationsURL(defaultQuery)
	request, err := http.NewRequest(method, fmt.Sprint(url), bytes.NewBuffer([]byte{}))
	if err != nil {
		log.Println("req", err)
		return map[int]string{}, nil
	}
	request.Header.Set("Content-type", "application/json")

	client := http.Client{
		Timeout: ha.Timeout,
	}
	resp, err := client.Do(request)
	if err != nil {
		log.Println("do error", err)
		return map[int]string{}, nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return map[int]string{}, nil
	}
	defer resp.Body.Close()

	var locations LocationsBody
	err = json.Unmarshal(body, &locations)
	if err != nil {
		log.Println("unmarshal error", err)
		return map[int]string{}, nil
	}

	m := make(map[int]string)
	for _, location := range locations.Locations {
		m[location.IntID] = location.ID
	}

	return m, nil
}
