package inmemcache

import (
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/patrickmn/go-cache"
)

// MemcacheClient is a useful helper and common interface for memcache client
type MemcacheClient interface {
	// Add according to bradfitz/gomemcache/memcache
	Add(item *memcache.Item) error
	// CompareAndSwap according to bradfitz/gomemcache/memcache
	CompareAndSwap(item *memcache.Item) error
	// Decrement according to bradfitz/gomemcache/memcache
	Decrement(key string, delta uint64) (newValue uint64, err error)
	// Delete according to bradfitz/gomemcache/memcache
	Delete(key string) error
	// DeleteAll according to bradfitz/gomemcache/memcache
	DeleteAll() error
	// FlushAll according to bradfitz/gomemcache/memcache
	FlushAll() error
	// Get according to bradfitz/gomemcache/memcache
	Get(key string) (item *memcache.Item, err error)
	// GetMulti according to bradfitz/gomemcache/memcache
	GetMulti(keys []string) (map[string]*memcache.Item, error)
	// Increment according to bradfitz/gomemcache/memcache
	Increment(key string, delta uint64) (newValue uint64, err error)
	// Replace according to bradfitz/gomemcache/memcache
	Replace(item *memcache.Item) error
	// Set according to bradfitz/gomemcache/memcache
	Set(item *memcache.Item) error
	// Touch according to bradfitz/gomemcache/memcache
	Touch(key string, seconds int32) (err error)
}

func newCache() *cache.Cache {
	return cache.New(cache.DefaultExpiration, cache.DefaultExpiration)
}

// New returns an inmemcache Client that implements Client
func New() *Client {
	return &Client{cache: newCache()}
}

// Client implements a local in-memory cache version of memcache client
type Client struct {
	MemcacheClient
	cache *cache.Cache
}

// Add according to bradfitz/gomemcache/memcache
func (c *Client) Add(item *memcache.Item) error {
	if _, ok := c.get(item.Key); !ok {
		return c.Set(item)
	}
	return memcache.ErrNotStored
}

// CompareAndSwap according to bradfitz/gomemcache/memcache
func (c *Client) CompareAndSwap(item *memcache.Item) error {
	// TODO: Implement me
	panic("not yet implemented")
}

// Decrement according to bradfitz/gomemcache/memcache
func (c *Client) Decrement(key string, delta uint64) (newValue uint64, err error) {
	// TODO: Implement me
	panic("not yet implemented")
}

// Delete according to bradfitz/gomemcache/memcache
func (c *Client) Delete(key string) error {
	if _, ok := c.get(key); ok {
		c.cache.Delete(key)
		return nil
	}
	return memcache.ErrCacheMiss
}

// DeleteAll according to bradfitz/gomemcache/memcache
func (c *Client) DeleteAll() error {
	for k := range c.cache.Items() {
		c.cache.Delete(k)
	}
	return nil
}

// FlushAll according to bradfitz/gomemcache/memcache
func (c *Client) FlushAll() error {
	c.cache.Flush()
	return nil
}

// Get according to bradfitz/gomemcache/memcache
func (c *Client) Get(key string) (item *memcache.Item, err error) {
	if item, found := c.get(key); found {
		return item, nil
	}
	return nil, memcache.ErrCacheMiss
}

// GetMulti according to bradfitz/gomemcache/memcache
func (c *Client) GetMulti(keys []string) (map[string]*memcache.Item, error) {
	m := make(map[string]*memcache.Item)
	for _, key := range keys {
		if item, found := c.get(key); found {
			m[item.Key] = item
		}
	}
	return nil, nil
}

// Increment according to bradfitz/gomemcache/memcache
func (c *Client) Increment(key string, delta uint64) (newValue uint64, err error) {
	// TODO: Implement me
	panic("not yet implemented")
}

// Replace according to bradfitz/gomemcache/memcache
func (c *Client) Replace(item *memcache.Item) error {
	if _, ok := c.get(item.Key); ok {
		return c.Set(item)
	}
	return memcache.ErrNotStored
}

// Set according to bradfitz/gomemcache/memcache
func (c *Client) Set(item *memcache.Item) error {
	c.cache.Set(item.Key, item, toDuration(item.Expiration))
	return nil
}

// Touch according to bradfitz/gomemcache/memcache
func (c *Client) Touch(key string, seconds int32) (err error) {
	if item, ok := c.get(key); ok {
		item.Expiration = seconds
		return c.Set(item)
	}
	return memcache.ErrCacheMiss
}

func (c *Client) get(key string) (item *memcache.Item, ok bool) {
	if i, ok := c.cache.Get(key); ok {
		return i.(*memcache.Item), true
	}
	return nil, false
}

const secondsPerMonth = 60 * 60 * 24 * 30

func toDuration(s int32) time.Duration {
	if s <= secondsPerMonth {
		return time.Second * time.Duration(s)
	}
	d := time.Now().Sub(time.Unix(int64(s), 0))
	if d > 0 {
		return d
	}
	return 0
}
