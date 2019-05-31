package cache

// 本地缓存 cache

import (
	"fmt"
	"sync"
	"time"
)

type Item struct {
	Value      interface{}
	Expiration int64
}

const (
	DefaultExpiration = 0
	NoExpiration      = -1
)

type Cache struct {
	defaultExpiration time.Duration
	bucket            map[string]Item
	mu                sync.RWMutex
	bucketSize        int
	queue             chan string
	stats             map[string]int64
}

func NewCache(defaultExpiration time.Duration) *Cache {
	cache := &Cache{
		defaultExpiration: defaultExpiration,
		bucket:            make(map[string]Item),
		queue:             make(chan string, 100),
		stats:             make(map[string]int64),
	}
	go cache.clean()
	return cache
}

func (c *Cache) set(k string, v interface{}, d time.Duration) {
	var e int64
	if d == DefaultExpiration {
		d = c.defaultExpiration
	}
	if d > 0 {
		e = time.Now().Add(d).UnixNano()
	}
	c.bucket[k] = Item{
		Value:      v,
		Expiration: e,
	}
	c.queue <- k
}

func (c *Cache) get(k string) (interface{}, bool) {
	item, found := c.bucket[k]
	if !found {
		return nil, false
	}
	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			return nil, false
		}
	}
	c.queue <- k
	return item.Value, true
}

func (c *Cache) clean() {
	visitTime := time.Now()
	ticker := time.NewTicker(time.Duration(10) * time.Second)
	for {
		select {
		case key := <-c.queue:
			c.stats[key]++
			visitTime = time.Now()
		case <-ticker.C:
			if time.Now().Sub(visitTime) > time.Duration(10*60)*time.Second {
				c.mu.Lock()
				for k, v := range c.bucket {
					if v.Expiration > 0 && time.Now().UnixNano() > v.Expiration {
						delete(c.bucket, k)
					}
				}
				c.mu.Unlock()
			}
		}
	}
}

func (c *Cache) Set(k string, v interface{}, d time.Duration) {
	var e int64
	if d == DefaultExpiration {
		d = c.defaultExpiration
	}
	if d > 0 {
		e = time.Now().Add(d).UnixNano()
	}
	c.mu.Lock()
	c.bucket[k] = Item{
		Value:      v,
		Expiration: e,
	}
	c.queue <- k
	c.mu.Unlock()
}

// Add an item to the cache, replacing any existing item, using the default
// expiration.
func (c *Cache) SetDefault(k string, v interface{}) {
	c.Set(k, v, DefaultExpiration)
}

// Add an item to the cache only if an item doesn't already exist for the given
// key, or if the existing item has expired. Returns an error otherwise.
func (c *Cache) Add(k string, v interface{}, d time.Duration) error {
	c.mu.Lock()
	_, found := c.get(k)
	if found {
		c.mu.Unlock()
		return fmt.Errorf("Item %s already exists", k)
	}
	c.set(k, v, d)
	c.mu.Unlock()
	return nil
}

// Set a new value for the cache key only if it already exists, and the existing
// item hasn't expired. Returns an error otherwise.
func (c *Cache) Replace(k string, v interface{}, d time.Duration) error {
	c.mu.Lock()
	_, found := c.get(k)
	if !found {
		c.mu.Unlock()
		return fmt.Errorf("Item %s doesn't exist", k)
	}
	c.set(k, v, d)
	c.mu.Unlock()
	return nil
}

// Get an item from the cache. Returns the item or nil, and a bool indicating
// whether the key was found.
func (c *Cache) Get(k string) (interface{}, bool) {
	c.mu.RLock()
	item, found := c.bucket[k]
	c.queue <- k
	if !found {
		c.mu.RUnlock()
		return nil, false
	}
	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			c.mu.RUnlock()
			return nil, false
		}
	}
	c.mu.RUnlock()
	return item.Value, true
}

func (c *Cache) GetAll() map[string]Item {
	c.mu.RLock()
	defer c.mu.RUnlock()
	m := make(map[string]Item, len(c.bucket))
	now := time.Now().UnixNano()
	for k, v := range c.bucket {
		if v.Expiration > 0 {
			if now > v.Expiration {
				continue
			}
		}
		m[k] = v
	}
	return m
}

// GetWithExpiration returns an item and its expiration time from the cache.
// It returns the item or nil, the expiration time if one is set (if the item
// never expires a zero value for time.Time is returned), and a bool indicating
// whether the key was found.
func (c *Cache) GetWithExpiration(k string) (interface{}, time.Time, bool) {
	c.mu.RLock()
	item, found := c.bucket[k]
	c.queue <- k
	if !found {
		c.mu.RUnlock()
		return nil, time.Time{}, false
	}

	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			c.mu.RUnlock()
			return nil, time.Time{}, false
		}

		// Return the item and the expiration time
		c.mu.RUnlock()
		return item.Value, time.Unix(0, item.Expiration), true
	}

	// If expiration <= 0 (i.e. no expiration time set) then return the item
	// and a zeroed time.Time
	c.mu.RUnlock()
	return item.Value, time.Time{}, true
}
