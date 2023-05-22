package main

import (
	"fmt"

	"github.com/dgraph-io/ristretto"
)

type Default struct {
	*ristretto.Cache
}

var cache, err = ristretto.NewCache(&ristretto.Config{
	NumCounters: 1e7,     // number of keys to track frequency of (10M).
	MaxCost:     1 << 30, // maximum cost of cache (1GB).
	BufferItems: 64,      // number of keys per Get buffer.
})

func DefaultNew() *Default {
	d := &Default{}
	if err != nil {
		panic(err)
	}
	d.Cache = cache

	if err != nil {
		panic(err)
	}
	return d
}

func main() {
	var defaultNew *Default
	for i := 0; i < 5; i++ {
		defaultNew = DefaultNew()
	}

	ok := defaultNew.Set("1", "2", 1)
	fmt.Println(ok)
	defaultNew.Wait()
	fmt.Println(defaultNew.Get("1"))
}

//func main() {
//	cache, err := ristretto.NewCache(&ristretto.Config{
//		NumCounters: 1e7,     // number of keys to track frequency of (10M).
//		MaxCost:     1 << 30, // maximum cost of cache (1GB).
//		BufferItems: 64,      // number of keys per Get buffer.
//	})
//	fmt.Println(err)
//	cache.Set("1", "2", 1)
//	cache.Wait()
//	fmt.Println(cache.Get("1"))
//}
