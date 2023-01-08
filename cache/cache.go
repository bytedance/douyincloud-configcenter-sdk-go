package cache

import (
	"sync"
)

type Cache struct {
	Config  map[string]*Item
	version string
	mu      sync.RWMutex
}

// NewCache creates instance of Cache
func NewCache() *Cache {
	return &Cache{Config: map[string]*Item{}}
}

type Item struct {
	Object interface{}
	Type   int64
}

func (c *Cache) Set(items map[string]*Item) {
	c.mu.Lock()
	newKeys := make(map[string]string, len(items))
	for k, v := range items {
		value := Item{Object: v.Object, Type: v.Type}
		c.Config[k] = &value
		newKeys[k] = ""
	}

	for k := range c.Config {
		if _, found := newKeys[k]; !found {
			delete(c.Config, k)
		}
	}

	c.mu.Unlock()
}

func (c *Cache) Get(k string) (*Item, bool) {
	c.mu.RLock()
	item, found := c.Config[k]
	if !found {
		c.mu.RUnlock()
		return nil, false
	}
	c.mu.RUnlock()
	return item, true
}

func (c *Cache) SetVersion(version string) {
	c.version = version
}

func (c *Cache) GetVersion() string {
	return c.version
}
