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
	cache.reapLoop()
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
	ticker := time.NewTicker(c.ttl / 10)

	go func() {
		for t := range ticker.C {
			for k, v := range c.values {
				elapsed := t.Sub(v.createdAt)
				if elapsed > c.ttl {
					delete(c.values, k)
				}
			}
		}
	}()

}
