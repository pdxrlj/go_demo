package main

import (
	"fmt"
	"time"

	"github.com/alitto/pond"
)

func main() {

	// Create an unbuffered (blocking) pool with a fixed
	// number of workers
	pool := pond.New(10, 5)

	// Submit 1000 tasks
	for i := 0; i < 1000; i++ {
		n := i
		pool.Submit(func() {
			fmt.Printf("Running task #%d\n", n)
			time.Sleep(time.Second * 3)
		})
	}

	// Stop the pool and wait for all submitted tasks to complete
	pool.StopAndWait()
}
