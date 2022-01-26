package cache

import (
	"fmt"
	"sync"
	"time"

	"github.com/asetriza/gocacheexpirable/fetch"
)

type Cacher interface {
	Get(id int) (string, error)
}

type Cache struct {
	C Cacher
}

func NewCacher(c Cacher) *Cache {
	return &Cache{
		C: c,
	}
}

type hcache struct {
	m      map[int]string
	mu     sync.RWMutex
	doneCh chan struct{}
	fetch  fetch.Fetcher
}

func New(f fetch.Fetcher) (*hcache, error) {
	hc := &hcache{
		doneCh: make(chan struct{}),
		mu:     sync.RWMutex{},
		fetch:  f,
	}

	if err := hc.fetchAll(); err != nil {
		return nil, err
	}

	return hc, nil
}

func (hc *hcache) fetchAll() error {
	hc.mu.Lock()
	defer hc.mu.Unlock()
	m, err := hc.fetch.FetchAll()
	if err != nil {
		return err
	}

	hc.m = m

	return nil
}

func (hc *hcache) Get(id int) (string, error) {
	hc.mu.RLock()
	value, found := hc.m[id]
	hc.mu.RUnlock()
	if !found {
		hc.mu.Lock()
		value, err := hc.fetch.Fetch(id)
		hc.mu.Unlock()
		if err != nil {
			return "", err
		}

		hc.mu.Lock()
		hc.m[id] = value
		hc.mu.Unlock()
		return value, err
	}

	return value, nil
}

func (hc *hcache) ReloadEvery(dur time.Duration) chan error {
	errCh := make(chan error)
	go hc.reloadIn(dur, errCh)
	return errCh
}

func (hc *hcache) reloadIn(dur time.Duration, errCh chan error) {
	t := time.NewTicker(dur)
	defer func() {
		t.Stop()
		close(hc.doneCh)
	}()
	for {
		select {
		case _ = <-hc.doneCh:
			return
		case _ = <-t.C:
			hc.clear()
			errCh <- hc.fetchAll()
		}
	}
}

func (hc *hcache) clear() {
	fmt.Println("Cache cleaned")
	hc.mu.Lock()
	defer hc.mu.Unlock()
	for key := range hc.m {
		delete(hc.m, key)
	}
}
