package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ejuju/gdb"
)

func main() {
	// init cache
	cache, err := gdb.NewCache(gdb.CacheConfig{
		ID:              "0",
		CleanupInterval: 500 * time.Millisecond,
		MaxItems:        500,
		SizeLimit:       1024,
	})
	if err != nil {
		// handle error ...
		log.Fatal(err)
		return
	}
	cache.Start()
	defer cache.Stop()

	err = cache.WriteOne(gdb.WriteOneRequest{
		ItemID: "0",
		Expiry: time.Now().Add(10 * time.Minute),
		Value:  []byte("hello"),
	})
	if err != nil {
		// handle error ...
		fmt.Println(err)
		return
	}
}
