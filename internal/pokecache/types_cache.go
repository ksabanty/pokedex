package pokecache

import (
    "time"
    "sync"
)

type cacheEntry struct {
    createdAt time.Time
    val []byte
}

type Cache struct {
    data map[string]cacheEntry
    mu *sync.Mutex
}

// NewCache() function that creates a new cache with a configurable interval (time.Duration)
func NewCache(interval time.Duration) Cache {
    c := Cache{
        data: make(map[string]cacheEntry),
        mu: &sync.Mutex{},
    }

    go c.reapLoop(interval)
    return c
}

func (c *Cache) Add(key string, val []byte) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.data[key] = cacheEntry{
        createdAt: time.Now().UTC(),
        val: val,
    }
}

func (c *Cache) Get(key string) ([]byte, bool) {
    c.mu.Lock()
    defer c.mu.Unlock()
    entry, ok := c.data[key]
    return entry.val, ok
}

func (c *Cache) reapLoop(interval time.Duration) {
    ticker := time.NewTicker(interval)
    for range ticker.C {
        c.reap(time.Now().UTC(), interval)
    }
}

func (c *Cache) reap(now time.Time, last time.Duration) {
    c.mu.Lock()
    defer c.mu.Unlock()
    for key, entry := range c.data {
        if entry.createdAt.Before(now.Add(-last)) {
            delete(c.data, key)
        }
    }
}
