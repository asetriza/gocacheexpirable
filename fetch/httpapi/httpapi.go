package httpapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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
	body, err := ha.makeRequest(method, url)
	if err != nil {
		return "", err
	}

	var location Locations
	err = json.Unmarshal(body, &location)
	if err != nil {
		log.Println("unmarshal error", err)
		return "", err
	}

	if len(location.Locations) == 0 {
		return "", errors.New("No location")
	}
	return fmt.Sprint(location.Locations[0].ID), nil
}

func (ha HttpAPI) FetchAll() (map[int]string, error) {

	method, url := ha.locationsURL(defaultQuery)
	body, err := ha.makeRequest(method, url)
	if err != nil {
		return map[int]string{}, err
	}

	var locations Locations
	err = json.Unmarshal(body, &locations)
	if err != nil {
		log.Println("unmarshal error", err)
		return map[int]string{}, err
	}

	m := make(map[int]string)
	for _, location := range locations.Locations {
		m[location.IntID] = location.ID
	}

	return m, nil
}

func (ha HttpAPI) makeRequest(method string, url *url.URL) ([]byte, error) {
	request, err := http.NewRequest(method, fmt.Sprint(url), bytes.NewBuffer([]byte{}))
	if err != nil {
		log.Println("req", err)
		return []byte{}, err
	}
	request.Header.Set("Content-type", "application/json")

	client := http.Client{
		Timeout: ha.Timeout,
	}
	resp, err := client.Do(request)
	if err != nil {
		log.Println("http client do error", err)
		return []byte{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return []byte{}, err
	}
	defer resp.Body.Close()

	return body, nil
}
