package gdb

import (
	"errors"
	"testing"
)

func TestWriteOne(t *testing.T) {
	// start cache
	c, err := NewCache(CacheConfig{
		id: "testuser",
	})
	if err != nil {
		t.Error(err)
		return
	}
	c.Start()
	defer c.Stop()

	// encode data
	encoded, err := encode("test data")
	if err != nil {
		t.Error(err)
		return
	}

	// write data
	err = c.writeOne(WriteOneRequest{
		itemID: "0",
		value:  encoded,
	})
	if err != nil {
		t.Error(err)
		return
	}

	// read data
	item, err := c.readOne(ReadOneRequest{
		itemID: "0",
	})
	str := ""
	err = item.DecodeInto(&str)
	if err != nil {
		t.Error(err)
		return
	}
	if str != "test data" {
		t.Error(err)
		return
	}

	// write large data
	largedata := make([]byte, c.config.sizelimit+1)
	err = c.writeOne(WriteOneRequest{
		itemID: "1",
		value:  largedata,
	})
	if err != nil {
		t.Error(err)
		return
	}

	// read large item
	item, err = c.readOne(ReadOneRequest{
		itemID: "1",
	})
	if err != nil {
		t.Error(err)
		return
	}
	if item.file == nil {
		t.Error("file should not be nil for large item")
		return
	}
}

func TestEraseOne(t *testing.T) {
	// start cache
	c, err := NewCache(CacheConfig{
		id: "testuser",
	})
	if err != nil {
		t.Error(err)
		return
	}
	c.Start()
	defer c.Stop()

	// encode data
	encoded, err := encode("test data")
	if err != nil {
		t.Error(err)
		return
	}

	// write data
	err = c.writeOne(WriteOneRequest{
		itemID: "0",
		value:  encoded,
	})
	if err != nil {
		t.Error(err)
		return
	}

	// delete data
	err = c.eraseOne(EraseOneRequest{
		itemID: "0",
	})
	if err != nil {
		t.Error(err)
		return
	}

	// check if item still exists
	_, err = c.readOne(ReadOneRequest{
		itemID: "0",
	})
	if !errors.Is(err, ErrUnknownID) {
		t.Error("unexpected result, ErrUnkownID was expected")
		return
	}
}
