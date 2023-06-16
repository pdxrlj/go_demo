package main

import (
	"fmt"
	"time"

	"golang.org/x/sync/singleflight"
)

var sg = singleflight.Group{}

func main() {
	for i := 0; i < 5; i++ {
		go func() {
			val, shared := sgf()
			fmt.Printf("out val: %v, shared: %v\n", val, shared)
		}()
	}

	val, shared := sgf()
	fmt.Printf("2 out val: %v, shared: %v\n", val, shared)

	select {}
}

func sgf() (int64, bool) {
	time.Sleep(time.Second * 3)
	val, err, shared := sg.Do("key", func() (interface{}, error) {
		return time.Now().UnixNano(), nil
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("int val: %v, shared: %v\n", val.(int64), shared)

	return val.(int64), shared
}
