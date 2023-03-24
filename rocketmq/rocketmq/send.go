package rocketmq

import (
	"context"
	"fmt"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

func Producer(tag string) {
	p, _ := rocketmq.NewProducer(
		producer.WithNameServer([]string{"192.168.1.65:9876"}), // 接入点地址
		producer.WithRetry(2), // 重试次数
		producer.WithGroupName("product_group"),
	)
	err := p.Start()
	if err != nil {
		panic(err)
	}
	defer func(p rocketmq.Producer) {
		_ = p.Shutdown()
	}(p)

	// 发送同步消息
	msg := &primitive.Message{
		Topic: "topic_test",
		Body:  []byte("Hello RocketMQ Go Client! " + tag),
	}
	msg.WithTag(tag)
	result, err := p.SendSync(context.Background(), msg)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Send sync message result: %s\n", result.String())
}
