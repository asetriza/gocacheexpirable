package cache

import "sync"

type hcache struct {
	m  map[int]string
	mu sync.RWMutex
}

func New() *hcache {
	return &hcache{
		m:  make(map[int]string),
		mu: sync.RWMutex{},
	}
}
