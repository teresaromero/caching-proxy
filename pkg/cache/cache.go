package cache

import (
	"container/list"
	"net/http"
	"sync"
	"time"
)

// Item represents the data stored in the cache.
type Item struct {
	ResponseBody       []byte
	ResponseHeaders    http.Header
	ResponseStatusCode int
	Expiration         time.Time
}

// Cache represents a cache with a fixed capacity and TTL.
type Cache struct {
	mu        sync.RWMutex
	itemsMap  map[string]*list.Element
	itemsList *list.List
	ttl       time.Duration
	capacity  int
	redis     *Redis
}

type CacheConfig struct {
	TTL      time.Duration
	Capacity int

	RedisAddr     string
	RedisDB       int
	RedisPwd      string
	RedisUsername string
}

// New creates a new Cache with the default capacity and TTL.
func New(config *CacheConfig) *Cache {
	return &Cache{
		itemsMap:  make(map[string]*list.Element),
		itemsList: list.New(),
		ttl:       config.TTL,
		capacity:  config.Capacity,
		redis:     NewRedis(config.RedisDB, config.RedisAddr, config.RedisUsername, config.RedisPwd, config.TTL),
	}
}

func (c *Cache) Get(key string) (*Item, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if element, ok := c.itemsMap[key]; ok {
		item, ok := element.Value.(*Item)
		if !ok {
			return nil, false
		}

		// implement TTL
		if item.Expiration.Before(time.Now()) {
			c.itemsList.Remove(element)
			delete(c.itemsMap, key)
			return nil, false
		}
		// implement LRU, used item should be moved to the front
		c.itemsList.MoveToFront(element)
		return item, true
	}
	return nil, false
}

func (c *Cache) Set(key string, item *Item) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// implement LRU, if the cache is full, remove the last item
	if c.itemsList.Len() >= c.capacity {
		c.itemsList.Remove(c.itemsList.Back())
	}

	// set item expiration with ttl
	item.Expiration = time.Now().Add(c.ttl)

	element := c.itemsList.PushFront(item)
	c.itemsMap[key] = element
}

func (c *Cache) RemoveAll() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.itemsMap = make(map[string]*list.Element)
	c.itemsList.Init()
}
