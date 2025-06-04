package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	cache map[string]cacheEntry
	mu    *sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		cache: make(map[string]cacheEntry),
		mu:    &sync.Mutex{},
	}
	go c.reapLoop(interval)
	return c
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) (val []byte, found bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	exists, ok := c.cache[key]
	if ok {
		return exists.val, true
	}
	return nil, false
}

func (c *Cache) reapLoop(reapInterval time.Duration) {
	ticker := time.NewTicker(5 * reapInterval)
	defer ticker.Stop()

	for {
		<-ticker.C
		currentTime := time.Now()

		c.mu.Lock()

		for key, i := range c.cache {
			if currentTime.Sub(i.createdAt) > reapInterval {
				delete(c.cache, key)
			}

		}
		c.mu.Unlock()
	}
}
