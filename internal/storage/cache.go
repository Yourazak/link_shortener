package storage

import (
	"errors"
	"sync"
)

type Cache struct {
	DB map[string]string
	mu sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		DB: make(map[string]string),
	}
}

func (c *Cache) Set(key, value string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, exists := c.DB[key]; !exists {
		c.DB[key] = value
		return nil
	}
	return errors.New("key is already exists")
}

func (c *Cache) Get(key string) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if value, exists := c.DB[key]; exists {
		return value, nil
	}
	return "", errors.New("key doest exists")
}
func (c *Cache) Close() error {
	return nil
}
