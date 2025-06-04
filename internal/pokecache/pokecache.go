package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	entry map[string]cacheEntry
	mutex sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	var cache Cache
	cache.entry = make(map[string]cacheEntry)

	go cache.reapLoop(interval)
	return &cache
}

func (c *Cache) Add(key string, val []byte) {
	c.mutex.Lock()
	var cache cacheEntry
	cache.val = val
	cache.createdAt = time.Now()
	c.entry[key] = cache
	c.mutex.Unlock()
}

func (c *Cache) Get(key string) (val []byte, b bool) {
	c.mutex.Lock()
	entry, ok := c.entry[key]
	defer c.mutex.Unlock()
	if ok {
		return entry.val, ok
	}
	return nil, ok
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		<-ticker.C
		c.mutex.Lock()
		for key := range c.entry {
			if interval < time.Since(c.entry[key].createdAt) {
				delete(c.entry, key)

			}
		}
		c.mutex.Unlock()
	}
}
