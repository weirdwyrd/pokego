package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	Entries      map[string]CacheEntry
	ReapInterval time.Duration
	mu           sync.Mutex
}

type CacheEntry struct {
	Timestamp time.Time
	EntryData []byte
}

func NewCache(reapInterval time.Duration) (*Cache, error) {
	cache := Cache{
		Entries:      make(map[string]CacheEntry),
		ReapInterval: reapInterval,
	}

	go cache.reapLoop()

	return &cache, nil
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.Entries[key]
	return entry.EntryData, ok
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Entries[key] = CacheEntry{
		Timestamp: time.Now(),
		EntryData: val,
	}
}

func (c *Cache) reapLoop() {
	for {
		c.mu.Lock()
		for k, v := range c.Entries {
			// remove the entry if it is older than the reap interval
			if time.Since(v.Timestamp) > c.ReapInterval {
				delete(c.Entries, k)
			}
		}
		c.mu.Unlock()
		time.Sleep(c.ReapInterval)
	}
}
