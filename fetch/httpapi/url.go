package httpapi

import (
	"net/url"
)

func (ha *HttpAPI) locationURL(ltype, key, value string) (string, *url.URL) {
	u := &url.URL{
		Scheme: ha.Scheme,
		Host:   ha.Host,
		Path:   ha.Endpoints.Location.Path,
	}

	rq := u.Query()

	rq.Set("type", ltype)
	rq.Set("key", key)
	rq.Set("value", value)

	u.RawQuery = rq.Encode()

	return ha.Endpoints.Location.Method, u
}

func (ha *HttpAPI) locationsURL(query string) (string, *url.URL) {
	u := &url.URL{
		Scheme: ha.Scheme,
		Host:   ha.Host,
		Path:   ha.Endpoints.Locations.Path,
	}

	rq := u.Query()

	rq.Set("query", query)

	u.RawQuery = rq.Encode()

	return ha.Endpoints.Locations.Method, u
}
