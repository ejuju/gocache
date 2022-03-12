package gocache

import (
	"errors"
	"testing"
)

func TestWriteOne(t *testing.T) {
	// start cache
	c, err := NewCache(CacheConfig{
		ID: "testuser",
	})
	if err != nil {
		t.Error(err)
		return
	}
	c.Start()
	defer c.Stop()

	// encode data
	testData := "test data"

	// write data
	err = c.WriteOne(WriteOneRequest{
		ItemID: "0",
		Value:  testData,
	})
	if err != nil {
		t.Error(err)
		return
	}

	// read data
	item, err := c.ReadOne(ReadOneRequest{
		ItemID: "0",
	})
	str := ""
	err = item.DecodeInto(&str)
	if err != nil {
		t.Error(err)
		return
	}
	if str != testData {
		t.Error("unexpected decoded result", "expected "+testData, "got "+str)
		return
	}

	// write large data
	largedata := make([]byte, c.config.SizeLimit+1)
	err = c.WriteOne(WriteOneRequest{
		ItemID: "1",
		Value:  largedata,
	})
	if err != nil {
		t.Error(err)
		return
	}

	// read large item
	item, err = c.ReadOne(ReadOneRequest{
		ItemID: "1",
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
		ID: "testuser",
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
	err = c.WriteOne(WriteOneRequest{
		ItemID: "0",
		Value:  encoded,
	})
	if err != nil {
		t.Error(err)
		return
	}

	// delete data
	err = c.EraseOne(EraseOneRequest{
		ItemID: "0",
	})
	if err != nil {
		t.Error(err)
		return
	}

	// check if item still exists
	_, err = c.ReadOne(ReadOneRequest{
		ItemID: "0",
	})
	if !errors.Is(err, ErrUnknownID) {
		t.Error("unexpected result, ErrUnkownID was expected")
		return
	}
}
