package gdb

import "errors"

// ErrMaxItemsReached occurs when the request fails because the maximum number of items in the map has been reached
// It is defined in the maxitems CacheConfig field
var ErrMaxItemsReached = errors.New("the maximum number of items has been reached")

// ErrUnknownID occurs when no resource was found for a given ID
var ErrUnknownID = errors.New("ID was not found")

// ErrEmptyItem occurs for attemps to decode items that don't have any data
var ErrEmptyItem = errors.New("item is empty, there's no data or file associated with it")
