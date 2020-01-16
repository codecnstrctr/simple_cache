package cache

import (
	"container/heap"
	"sync"
	"time"
)

// Cache represents concurrent-safe registry
// for key-value pairs with ttl
type Cache struct {
	sync.RWMutex
	storage   map[string]interface{}
	ttlHeap   *timeoutHeap
	ttlTicker *time.Ticker
	done      chan struct{}
}

// NewCache creates new cache instance
func NewCache(removeExpiredPeriod int) *Cache {
	cache := new(Cache)

	cache.storage = make(map[string]interface{})
	cache.ttlHeap = newTimeoutHeap()
	cache.ttlTicker = time.NewTicker(time.Second * time.Duration(removeExpiredPeriod))
	cache.done = make(chan struct{})

	go func(c *Cache) {
		select {
		case <-c.ttlTicker.C:
			c.removeExpired()
		case <-c.done:
			break
		}
	}(cache)

	return cache
}

// Stop stops c self-cleaning
func (c *Cache) Stop() {
	c.done <- struct{}{}
	close(c.done)
	c.ttlTicker.Stop()
}

// Set sets the key-data pair in c
func (c *Cache) Set(key string, data interface{}, ttl time.Duration) {
	c.Lock()
	c.set(key, data, ttl)
	c.Unlock()
}

// Get returns the data stored at key from c
func (c *Cache) Get(key string) (interface{}, bool) {
	c.RLock()
	data, ok := c.get(key)
	c.RUnlock()

	return data, ok
}

// Remove deletes the data stored at key from c
func (c *Cache) Remove(key string) bool {
	c.Lock()
	ok := c.remove(key)
	c.Unlock()

	return ok
}

// Keys returns the list of keys stored in c
func (c *Cache) Keys() []string {
	c.RLock()
	keys := c.keys()
	c.RUnlock()

	return keys
}

func (c *Cache) set(key string, data interface{}, ttl time.Duration) {
	c.storage[key] = data
	if ttl == 0 {
		return
	}

	heap.Push(c.ttlHeap,
		timeout{
			expireAt: time.Now().Add(ttl),
			key:      key,
		})
}

func (c *Cache) get(key string) (interface{}, bool) {
	v, ok := c.storage[key]
	return v, ok
}

func (c *Cache) remove(key string) bool {
	_, ok := c.get(key)
	if ok {
		delete(c.storage, key)
	}

	return ok
}

func (c *Cache) keys() []string {
	keys := make([]string, 0)
	for k := range c.storage {
		keys = append(keys, k)
	}

	return keys
}

func (c *Cache) removeExpired() {
	c.Lock()
	for {
		top, ok := c.ttlHeap.take().(timeout)
		if !ok || top.expireAt.After(time.Now()) {
			break
		}

		heap.Pop(c.ttlHeap)
		c.remove(top.key)
	}
	c.Unlock()
}
