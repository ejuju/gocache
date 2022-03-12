package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ejuju/gocache"
)

func main() {
	// init cache
	cache, err := gocache.NewCache(gocache.CacheConfig{
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

	// write one item
	err = cache.WriteOne(gocache.WriteOneRequest{
		ItemID: "0",
		Expiry: time.Now().Add(10 * time.Minute),
		Value:  "hello",
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	// get one item
	item, err := cache.ReadOne(gocache.ReadOneRequest{
		ItemID: "0",
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	// decode item into var
	intovar := ""
	err = item.DecodeInto(&intovar)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(intovar) // > "hello"
}
