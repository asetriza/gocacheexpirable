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

type hcache struct {
	m      map[int]string
	mu     sync.RWMutex
	doneCh chan struct{}
	fetch  fetch.Fetcher
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
	hc.doneCh = make(chan struct{})
	return hc, nil
}

func (hc *hcache) Get(id int) (string, error) {
	hc.mu.RLock()
	value, found := hc.m[id]
	hc.mu.RUnlock()
	if !found {
		value, err := hc.fetch.Fetch(id)
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

func (hc *hcache) ReloadEvery(dur time.Duration) {
	go hc.reloadIn(dur)
}

func (hc *hcache) reloadIn(dur time.Duration) {
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
		}
	}
}

func (hc *hcache) clear() {
	fmt.Println("reload started")
	hc.mu.Lock()
	defer hc.mu.Unlock()
	for key := range hc.m {
		delete(hc.m, key)
	}
}
