package main

import (
	"context"
	"fmt"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

func main() {
	p, _ := rocketmq.NewProducer(
		producer.WithNameServer([]string{"192.168.1.65:9876"}), // 接入点地址
		producer.WithRetry(2),                  // 重试次数
		producer.WithGroupName("ProductGroup"), // 分组名称
	)
	err := p.Start()
	if err != nil {
		panic(err)
	}
	defer func(p rocketmq.Producer) {
		_ = p.Shutdown()
	}(p)

	// 发送同步消息
	result, err := p.SendSync(context.Background(), &primitive.Message{
		Topic: "TopicTest",
		Body:  []byte("Hello RocketMQ Go Client!"),
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Send sync message result: %s", result.String())
}
