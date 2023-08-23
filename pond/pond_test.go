package main

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/alitto/pond"
)

func TestDynamic_Size(t *testing.T) {
	// Create a buffered (non-blocking) pool that can scale up to 100 workers
	// and has a buffer capacity of 1000 tasks
	pool := pond.New(3, 200)

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

func TestFixed_Size(t *testing.T) {
	// Create an unbuffered (blocking) pool with a fixed
	// number of workers
	pool := pond.New(3, 0, pond.MinWorkers(10))

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

func TestGroup_Tasks(t *testing.T) {
	// Create a pool
	pool := pond.New(3, 1000)
	defer pool.StopAndWait()

	// Create a task group
	group := pool.Group()

	// Submit a group of tasks
	for i := 0; i < 20; i++ {
		n := i
		group.Submit(func() {
			fmt.Printf("Running group task #%d\n", n)
			time.Sleep(time.Second * 3)
		})
	}

	// Wait for all tasks in the group to complete
	group.Wait()
}

func TestGroup_Context(t *testing.T) {
	//debug.SetMemoryLimit()
	// Create a worker pool
	pool := pond.New(2, 1000)
	defer pool.StopAndWait()

	// Create a task group associated to a context
	group, ctx := pool.GroupContext(context.Background())

	var urls = []string{
		"https://www.golang.org/",
		"https://www.google.com/",
		"https://www.github.com/",
	}

	// Submit tasks to fetch each URL
	for _, url := range urls {
		url := url
		group.Submit(func() error {
			fmt.Println("Fetching URL:", url)
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
			resp, err := http.DefaultClient.Do(req)
			if err == nil {
				resp.Body.Close()
			}
			return err
		})
	}

	// Wait for all HTTP requests to complete.
	err := group.Wait()
	if err != nil {
		fmt.Printf("Failed to fetch URLs: %v", err)
	} else {
		fmt.Println("Successfully fetched all URLs")
	}
}
