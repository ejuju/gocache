package gdb

import (
	"os"
	"sync"
	"time"
)

// Cache provides in-memory data storage
// Use the NewCache func to create a new cache
type Cache struct {
	mu      sync.RWMutex     // rwmutex to avoid data races
	items   map[string]*Item // actual data
	stop    chan struct{}    // stop chan
	config  CacheConfig      // configuration (settings/limits/constraints)
	dirpath string           // local directory path
}

// CacheConfig holds the cache configuration
// id sets the cache id (used for error tracing and local file paths related to the cache) (default: time.Now())
// interval sets the cleanup interval for the cleanup routine (default: one second)
// maxitems sets the maximum number of items (default: one hundred thousand)
// sizelimit sets the number of bytes for an item to be considered too big for in-memory (default: 500 kilobyte)
// an item with a size that exceeds the sizelimit will be stored a to a file
type CacheConfig struct {
	ID              string
	MaxItems        int
	SizeLimit       int
	CleanupInterval time.Duration
}

// NewCache instantiates a cache
// it calls make(map[string]*item to avoid assigning to nil map for the cache.items field)
// it also creates a local directory to store temporary files
func NewCache(config CacheConfig) (*Cache, error) {
	// validate config, set defaults if needed
	if config.ID == "" {
		config.ID = time.Now().String()
	}
	if config.CleanupInterval <= 0 {
		config.CleanupInterval = 1 * time.Second
	}
	if config.MaxItems <= 0 {
		config.MaxItems = 100 * 1000
	}
	if config.SizeLimit <= 0 {
		config.SizeLimit = 500 * 1024
	}

	// define local path
	dirpath := "_gdb_cache_" + config.ID

	// remove old directory and files
	err := os.RemoveAll(dirpath)
	if err != nil {
		return nil, err
	}

	// create dir
	err = os.Mkdir(dirpath, os.ModePerm)
	if err != nil {
		return nil, err
	}

	// init cache and return
	return &Cache{
		items:   make(map[string]*Item),
		stop:    make(chan struct{}),
		dirpath: dirpath,
		config:  config,
	}, nil
}

// Start prepares the cache for being used
// it starts the cleanup loop that will take care of removing expired items
func (c *Cache) Start() {
	go c.startCleanupLoop(c.config.CleanupInterval)
}

// Stop gracefully stops the cache
func (c *Cache) Stop() {
	close(c.stop)
}

// removeExpiredItems iterates over the items map and deletes expired items
func (c *Cache) removeExpiredItems() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, item := range c.items {
		if item.isExpired() {
			delete(c.items, item.id)
		}
	}
}

// startCleanupLoop removes expired items at the given interval until the cache's stop channel gets closed
func (c *Cache) startCleanupLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-c.stop:
			return
		case <-ticker.C:
			c.removeExpiredItems()
		}
	}
}
