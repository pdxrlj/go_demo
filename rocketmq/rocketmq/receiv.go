package main

import (
	"context"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

func main() {
	pullConsumer, err := rocketmq.NewPushConsumer(consumer.WithNameServer([]string{"192.168.1.65:9876"}), // 接入点地址
		consumer.WithConsumerModel(consumer.Clustering),
		consumer.WithGroupName("ConsumerGroup"), // 分组名称
	)
	if err != nil {
		panic(err)
	}
	err = pullConsumer.Start()
	if err != nil {
		panic(err)
	}
	defer func(pullConsumer rocketmq.PushConsumer) {
		err := pullConsumer.Shutdown()
		if err != nil {
			panic(err)
		}
	}(pullConsumer)

	err = pullConsumer.Subscribe("TopicTest", consumer.MessageSelector{}, func(ctx context.Context, ext ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for _, msg := range ext {
			println(msg.Body)
		}
		return consumer.ConsumeSuccess, nil
	})
	if err != nil {
		panic(err)
	}

}
