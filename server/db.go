package server

import (
	"sync"
)

type memoryDB struct {
	items map[string]string
	mu    *sync.RWMutex
}

func newDB() memoryDB {
	return memoryDB{
		items: map[string]string{},
	}
}

func (m memoryDB) set(key, value string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.items[key] = value
}

func (m memoryDB) get(key string) (string, bool) {
	m.mu.RLock()
	defer m.mu.Unlock()
	val, found := m.items[key]
	return val, found
}

func (m memoryDB) delete(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.items, key)
}
