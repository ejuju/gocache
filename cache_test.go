package gdb

import "testing"

func TestNewCache(t *testing.T) {
	c, err := NewCache(CacheConfig{
		id: "testuser",
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
