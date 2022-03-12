package gocache

import (
	"errors"
	"testing"
	"time"
)

func TestNewCache(t *testing.T) {
	c, err := NewCache(CacheConfig{
		ID: "testuser",
	})
	if err != nil {
		t.Error(err)
		return
	}
	if c.items == nil {
		t.Error("items map should not be nil")
		return
	}
}

func TestCacheStart(t *testing.T) {
	c, err := NewCache(CacheConfig{
		ID:              "testuser",
		CleanupInterval: 1 * time.Millisecond,
	})
	if err != nil {
		t.Error(err)
		return
	}
	c.Start()
	defer c.Stop()

	// check if cleanup look is working
	err = c.WriteOne(WriteOneRequest{
		ItemID: "0",
		Value:  []byte("test"),
		Expiry: time.Now().Add(1 * time.Millisecond),
	})
	if err != nil {
		t.Error(err)
		return
	}

	// read data before expiry
	_, err = c.ReadOne(ReadOneRequest{
		ItemID: "0",
	})
	if err != nil {
		t.Error(err)
		return
	}

	// read data after expiry
	time.Sleep(1 * time.Millisecond)
	_, err = c.ReadOne(ReadOneRequest{
		ItemID: "0",
	})
	if !errors.Is(err, ErrUnknownID) {
		t.Error("item should have been delete by cleanup routine")
		return
	}
}
