package cache

import (
	"sync"
	"time"
)

type CacheItem struct {
	Data      interface{}
	Timestamp time.Time
}

type Cache struct {
	Data  map[string]CacheItem
	TTL   time.Duration
	Mutex sync.RWMutex
}

func NewCache(ttl time.Duration) *Cache {
	return &Cache{
		Data: make(map[string]CacheItem),
		TTL:  ttl,
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.Mutex.RLock()
	defer c.Mutex.RUnlock()

	item, exists := c.Data[key]
	if !exists || time.Since(item.Timestamp) > c.TTL {
		return nil, false
	}
	return item.Data, true
}

func (c *Cache) Set(key string, data interface{}) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	c.Data[key] = CacheItem{
		Data:      data,
		Timestamp: time.Now(),
	}
}

func (c *Cache) Delete(key string) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	delete(c.Data, key)
}
