package fetch

type Fetcher interface {
	// Fetch
	Fetch(id int) (string, error)

	// FetchAll
	FetchAll() (map[int]string, error)
}

type Fetch struct {
	F Fetcher
}

func New(f Fetcher) Fetch {
	return Fetch{
		F: f,
	}
}
