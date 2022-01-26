package cache

import (
	"sync"

	"github.com/asetriza/gocacheexpirable/fetch"
)

type Cacher interface {
	Get(id int) (string, error)
}

type hcache struct {
	m     map[int]string
	mu    sync.RWMutex
	fetch fetch.Fetcher
}

func New(f fetch.Fetcher) (*hcache, error) {
	hc := &hcache{}
	hc.fetch = f
	m, err := hc.fetch.FetchAll()
	if err != nil {
		return nil, err
	}
	hc.m = m
	hc.mu = sync.RWMutex{}
	return hc, nil
}

func (hc *hcache) Get(id int) (string, error) {
	hc.mu.RLock()
	value, found := hc.m[id]
	if !found {
		value, err := hc.fetch.Fetch(id)
		if err != nil {
			return "", err
		}
		hc.mu.Lock()
		hc.m[id] = value
		hc.mu.Unlock()
		hc.mu.RUnlock()
		return value, err
	}
	hc.mu.RUnlock()
	return value, nil
}
