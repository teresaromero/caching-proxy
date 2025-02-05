package cache

import (
	"container/list"
	"context"
	"log"
	"net/http"
	"sync"
	"time"
)

// Item represents the data stored in the cache.
type Item struct {
	Key                string
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

func (c *Cache) Get(ctx context.Context, key string) (*Item, bool) {
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

	// if no redis is configured, return here
	if c.redis == nil {
		return nil, false
	}

	// if the item is not in the cache, check if it is in the redis
	log.Println("Looking for item in redis...")
	if item, err := c.redis.Get(ctx, key); err != nil {
		item, ok := item.(*Item)
		if !ok {
			return nil, false
		}

		// set item in-memory cache to avoid multiple redis calls
		c.Set(key, item)
		return item, true
	}
	return nil, false
}

func (c *Cache) Set(key string, item *Item) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// implement LRU, if the cache is full, remove the last item
	if c.itemsList.Len() >= c.capacity {
		back := c.itemsList.Back()
		c.itemsList.Remove(back)
		delete(c.itemsMap, back.Value.(*Item).Key)
	}

	element := c.itemsList.PushFront(item)
	c.itemsMap[key] = element

	if c.redis != nil {
		// set item in redis asynchronously
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := c.redis.Set(ctx, key, item); err != nil {
				log.Println("Error setting item to redis:", err)
			}
		}()
	}
}

func (c *Cache) RemoveAll(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.itemsMap = make(map[string]*list.Element)
	c.itemsList.Init()

	if c.redis != nil {
		return c.redis.RemoveAll(ctx)
	}
	return nil
}

func (c *Cache) TTL() time.Duration {
	return c.ttl
}
