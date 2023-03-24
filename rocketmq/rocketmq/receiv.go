package rocketmq

import (
	"context"
	"fmt"
	"os"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

func Receiver(tag string) {
	sig := make(chan os.Signal)

	pullConsumer, err := rocketmq.NewPushConsumer(consumer.WithNameServer([]string{"192.168.1.65:9876"}), // 接入点地址
		consumer.WithConsumerModel(consumer.Clustering),
		consumer.WithGroupName("consumer_group_"+tag), // 分组名称
		consumer.WithAutoCommit(true),
		consumer.WithConsumerOrder(true),
		consumer.WithConsumeFromWhere(consumer.ConsumeFromFirstOffset),
	)
	if err != nil {
		panic(err)
	}

	err = pullConsumer.Subscribe("topic_test", consumer.MessageSelector{
		Type:       consumer.TAG,
		Expression: tag,
	}, func(ctx context.Context, ext ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for _, msg := range ext {
			fmt.Printf("=========================== receiver tag:%s msg:%s\n", tag, string(msg.Body))
		}
		return consumer.ConsumeSuccess, nil
	})
	if err != nil {
		panic(err)
	}

	err = pullConsumer.Start()
	if err != nil {
		panic(err)
	}
	<-sig
	err = pullConsumer.Shutdown()
	if err != nil {
		panic(err)
	}

}
