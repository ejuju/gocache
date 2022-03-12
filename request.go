package gocache

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
	ItemID     string
	Value      interface{}
	Expiry     time.Time
}

//
type EraseOneRequest struct {
	receivedAt time.Time
	ItemID     string
}

//
type ReadOneRequest struct {
	receivedAt time.Time
	ItemID     string
	// filters []Filter
}

// //
// type Filter struct {
// }

// WriteOne creates or update an item if the request is valid
// when the item exceeds the size limit, it is stored to a temporary file
func (c *Cache) WriteOne(woreq WriteOneRequest) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// check max items
	if len(c.items) >= c.config.MaxItems {
		return ErrMaxItemsReached
	}

	// create new item struct
	item := &Item{
		id:     woreq.ItemID,
		expiry: woreq.Expiry,
	}

	// encode data
	encoded, err := encode(woreq.Value)
	if err != nil {
		return err
	}

	// if max item size exceeded, write data to temp file and store pointer to file in map
	if len(encoded) > c.config.SizeLimit {
		filename := "_bigitem_*.gdbi"
		f, err := os.CreateTemp(c.dirpath, filename)
		if err != nil {
			return err
		}
		defer f.Close()

		// store data to file and set pointer to file on item
		_, err = f.Write(encoded)
		if err != nil {
			return err
		}
		item.file = f
		c.items[woreq.ItemID] = item
		return nil
	}

	// add encoded data to item struct and write item to map
	item.data = encoded
	c.items[woreq.ItemID] = item

	return nil
}

// EraseOne deletes an item if the request is valid
func (c *Cache) EraseOne(eoreq EraseOneRequest) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, exists := c.items[eoreq.ItemID]
	if !exists {
		return ErrUnknownID
	}

	delete(c.items, eoreq.ItemID)
	return nil
}

// ReadOne reads an item if the request is valid
func (c *Cache) ReadOne(roreq ReadOneRequest) (*Item, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.items[roreq.ItemID]
	if !ok {
		return nil, ErrUnknownID
	}

	return item, nil
}
