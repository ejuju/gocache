package gdb

import (
	"os"
	"time"
)

// store requests in file
// following fmt:
// {opType uint8}+":"+{operation writeOperation | deleteOperation etc.}

// Keys to identify request types
const (
	KeyWriteOne uint8 = iota
	KeyEraseOne uint8 = iota
	KeyReadOne  uint8 = iota
)

// WriteOneRequest represents a request from the client to add or update an item in the cache
// Write operations don't check if the item already exists:
// so if you want to make sure an item can't be overwritten, check if the item exists before
type WriteOneRequest struct {
	receivedAt time.Time
	itemID     string
	value      []byte
	expiry     time.Time
}

//
type EraseOneRequest struct {
	receivedAt time.Time
	itemID     string
}

//
type ReadOneRequest struct {
	receivedAt time.Time
	itemID     string
	// filters []Filter
}

// //
// type Filter struct {
// }

// writeOne creates or update an item if the request is valid
func (c *Cache) writeOne(woreq WriteOneRequest) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// check max items
	if len(c.items) >= c.config.maxitems {
		return ErrMaxItemsReached
	}

	// create new item struct
	item := &Item{
		id:     woreq.itemID,
		expiry: woreq.expiry,
	}

	// if max item size exceeded, write data to file and store pointer to file in map
	if len(woreq.value) > c.config.sizelimit {
		filename := "_bigitem_*.gdbi"
		f, err := os.CreateTemp(c.dirpath, filename)
		if err != nil {
			return err
		}
		for {
		}
		// encode data
		encoded, err := encode(woreq.value)
		if err != nil {
			return err
		}
		// store data to file and set pointer to file on item
		_, err = f.Write(encoded)
		if err != nil {
			return err
		}
		item.file = f
		c.items[woreq.itemID] = item
		return nil
	}

	// add data to item struct and write item to map
	item.data = woreq.value
	c.items[woreq.itemID] = item

	return nil
}

// eraseOne deletes an item if the request is valid
func (c *Cache) eraseOne(eoreq EraseOneRequest) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, exists := c.items[eoreq.itemID]
	if !exists {
		return ErrUnknownID
	}

	delete(c.items, eoreq.itemID)
	return nil
}

// readOne reads an item if the request is valid
func (c *Cache) readOne(roreq ReadOneRequest) (*Item, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.items[roreq.itemID]
	if !ok {
		return nil, ErrUnknownID
	}

	return item, nil
}
