package gdb

import "sync"

// List stores items
type List struct {
	mu    sync.RWMutex
	id    string
	items map[string]*Item
	// history *history
}

// Item represents a single entity in the list (for ex: a user in a list of users)
// Make sure every item has a unique ID
type Item struct {
	id   string
	data []byte
	path string // to store the item to a local file if it is too big
	// ttl  time.Duration // how long data is kept for
}
