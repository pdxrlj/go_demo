package main

import (
	"time"

	"github.com/apache/rocketmq-client-go/v2/rlog"

	"rocketmq/admin"
	"rocketmq/rocketmq"
)

func main() {
	admin.CreateTopic()
	rlog.SetLogLevel("error")
	go func() {
		tags := []string{"tag1", "tag2", "tag3"}
		for i := 0; i < len(tags); i++ {
			for y := 0; y < 3; y++ {
				rocketmq.Producer(tags[i])
				time.Sleep(time.Second * 1)
			}
		}

	}()

	go func() {
		time.Sleep(time.Second * 5)
		tags := []string{"tag1", "tag2", "tag3"}
		for _, tag := range tags {
			go func(t string) {
				rocketmq.Receiver(t)
			}(tag)
		}

	}()

	select {}
}
