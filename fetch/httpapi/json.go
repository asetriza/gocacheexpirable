package httpapi

type LocationBody struct {
	Locations []struct {
		IntID int `json:"int_id"`
	} `json:"locations"`
}

type LocationsBody struct {
	Locations []Location `json:"locations"`
}

type Location struct {
	IntID int    `json:"int_id"`
	ID    string `json:"id"`
}
