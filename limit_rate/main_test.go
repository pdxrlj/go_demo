package main

import (
	"testing"
	"time"
)

func Test_Request(t *testing.T) {
	for i := 0; i < 10; i++ {
		go TestRequest("token_1")
	}
	time.Sleep(time.Second * 2)
}
