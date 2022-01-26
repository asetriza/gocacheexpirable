package httpapi

type Locations struct {
	Locations []Location `json:"locations"`
}

type Location struct {
	IntID int    `json:"int_id"`
	ID    string `json:"id"`
}
