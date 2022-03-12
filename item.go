package gdb

import (
	"io/ioutil"
	"os"
	"time"
)

// Item represents a single entity in the list (for ex: a user in a list of users)
// make sure every item has a unique ID
type Item struct {
	id     string
	data   []byte
	file   *os.File  // to store the item to a local file if it is too big
	expiry time.Time // how long data is kept for
}

func (i *Item) isExpired() bool {
	// expiry time not set, not expired
	if i.expiry.IsZero() {
		return false
	}
	// item is expired
	if time.Now().After(i.expiry) {
		return true
	}
	// not expired
	return false
}

// DecodeInto decodes the item in the provided variable (must be a pointer)
func (i *Item) DecodeInto(intoptr interface{}) error {
	// item has in-memory data
	if len(i.data) > 0 {
		return decode(i.data, intoptr)
	}

	// item has data on file
	if i.file != nil {
		fdata, err := ioutil.ReadAll(i.file)
		if err != nil {
			return err
		}
		return decode(fdata, intoptr)
	}

	// item is empty
	return ErrEmptyItem
}
