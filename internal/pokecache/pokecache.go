package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	values map[string]*cacheEntry
	mu     *sync.RWMutex
	ttl    time.Duration
}

func NewCache(interval time.Duration) *Cache {
	cache := Cache{
		values: map[string]*cacheEntry{},
		mu:     &sync.RWMutex{},
		ttl:    interval,
	}
	go cache.reapLoop()
	return &cache
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry := cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
	c.values[key] = &entry
}

func (c *Cache) Get(key string) (val []byte, ok bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.values[key]
	if !ok {
		return nil, ok
	}
	return entry.val, ok
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.ttl)
	for range ticker.C {
		reap(c)
	}
}

func reap(c *Cache) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for k, v := range c.values {
		minBirthTime := time.Now().Add(-c.ttl)
		if v.createdAt.Before(minBirthTime) {
			delete(c.values, k)
		}
	}
}
