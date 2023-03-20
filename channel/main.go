package main

import (
	"fmt"
	"time"
)

func main() {
	TestSendClose()
}

// TestSendClose 测试发送后关闭
func TestSendClose() {
	ch := make(chan int, 3)
	go func() {
		for i := 0; i < 5; i++ {
			ch <- i
		}
		close(ch)
	}()

	go func() {
		for item := range ch {
			fmt.Println(item)
			time.Sleep(time.Second * 1)
		}
	}()

	time.Sleep(time.Second * 10)
}
