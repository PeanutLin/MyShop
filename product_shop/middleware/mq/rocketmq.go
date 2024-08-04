package mq

import (
	"context"
	"encoding/json"
	"productshop/product_shop/middleware/logs"
	"productshop/product_shop/middleware/mq_content"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/pkg/errors"
)

var (
	OrderProducer rocketmq.Producer // 用于实现异步订单落库
)

// GROUP
var (
	OrderGroup = "ORDER_GROUP"
)

// TAG
var (
	OrderTag = "ORDER_TAG"
)

// TOPIC
var (
	OrdetTopic = "ORDER_TOPIC"
)

// KEY
var (
	OrderKey = "ORDER_KEY"
)

type MQfunc func(ctx context.Context, msg ...*primitive.MessageExt) (consumer.ConsumeResult, error)

func MustInitBookProductProducer() {
	var err error
	OrderProducer, err = rocketmq.NewProducer(
		producer.WithNameServer([]string{"172.18.105.29:9876"}), // 接入点地址
		producer.WithRetry(2),              // 重试次数
		producer.WithGroupName(OrderGroup), // 分组名称
	)
	if err != nil {
		panic(err)
	}

	OrderProducer.Start()
}

func SendOrderMessage(ctx context.Context, content *mq_content.Message) error {
	contentJSON, err := json.Marshal(content)
	if err != nil {
		logs.Error("[json.Marshal] error", logs.String("error msg", err.Error()))
		return errors.Wrap(err, "json error")
	}
	msg := &primitive.Message{
		Topic: OrdetTopic,
		Body:  contentJSON,
	}
	msg.WithTag(OrderTag)
	msg.WithKeys([]string{OrderKey})
	_, err = OrderProducer.SendSync(ctx, msg)
	if err != nil {
		return err
	}
	return nil
}

func MustStartPushConsumeBookProductMQ(fn MQfunc) {
	c, err := rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{"172.18.105.29:9876"}), // 接入点地址
		consumer.WithConsumerModel(consumer.Clustering),
		consumer.WithGroupName(OrderGroup), // 分组名称
	)
	if err != nil {
		panic(err)
	}

	selector := consumer.MessageSelector{
		Type:       consumer.TAG,
		Expression: OrderTag,
	}

	err = c.Subscribe(OrdetTopic, selector, fn)
	if err != nil {
		panic(err)
	}

	forever := make(chan struct{})

	err = c.Start()
	if err != nil {
		panic(err)
	}

	<-forever
}

func MustStartPullConsumeBookProductMQ() {
	c, err := rocketmq.NewPullConsumer(
		consumer.WithNameServer([]string{"172.18.105.29:9876"}), // 接入点地址
		consumer.WithConsumerModel(consumer.BroadCasting),
		consumer.WithGroupName(OrderGroup), // 分组名称
	)
	if err != nil {
		panic(err)
	}

	selector := consumer.MessageSelector{
		Type:       consumer.TAG,
		Expression: OrderTag,
	}

	err = c.Subscribe(OrdetTopic, selector)
	if err != nil {
		panic(err)
	}

	err = c.Start()
	if err != nil {
		panic(err)
	}

	panic("todo")
}
