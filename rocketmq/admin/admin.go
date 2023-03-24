package admin

import (
	"context"
	"fmt"

	"github.com/apache/rocketmq-client-go/v2/admin"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

func Client() admin.Admin {
	resolver := primitive.NewPassthroughResolver([]string{"192.168.1.65:9876"})
	newAdmin, err := admin.NewAdmin(admin.WithResolver(resolver))
	if err != nil {
		panic(err)
	}
	return newAdmin
}

func CreateTopic() {
	err := Client().CreateTopic(
		context.Background(),
		admin.WithTopicCreate("topic_test"),
		admin.WithBrokerAddrCreate("192.168.1.65:10911"),
	)
	if err != nil {
		fmt.Println("Create topic error:", err.Error())
	}
}
