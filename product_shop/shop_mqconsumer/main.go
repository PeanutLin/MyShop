package main

import (
	"context"
	"encoding/json"
	"productshop/product_shop/middleware/logs"
	"productshop/product_shop/middleware/mq"
	"productshop/product_shop/middleware/mq_content"
	"productshop/product_shop/middleware/mysql"

	repository "productshop/product_shop/base_repository"
	service "productshop/product_shop/base_service"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

var (
	productService service.IProductService
	orderService   service.IOrderService
)

func OrderConsumer(ctx context.Context, msg ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	for _, m := range msg {
		message := &mq_content.Message{}
		err := json.Unmarshal([]byte(m.Body), message)
		if err != nil {
			logs.Error("[json Unmarshal] error", logs.String("msg err", err.Error()))
			return consumer.ConsumeRetryLater, err
		}
		// 插入订单
		_, err = orderService.InsertOrderByMessage(ctx, message)
		if err != nil {
			logs.Error(err.Error())
			return consumer.ConsumeRetryLater, err
		}

		// 扣除商品数量
		err = productService.SubNumber(ctx, message.ProductID, message.ProductNum)
		if err != nil {
			logs.Error(err.Error())
			return consumer.ConsumeRetryLater, err
		}
	}
	return consumer.ConsumeSuccess, nil
}

func main() {
	defer logs.Init().Sync()
	mysql.Init()

	// product 服务
	productRepository := repository.NewProductRepository()
	productService = service.NewProductService(productRepository)

	// order 服务
	orderRepository := repository.NewOrderRepository()
	orderService = service.NewOrderService(orderRepository)

	// 启动消费者服务
	mq.MustStartPushConsumeBookProductMQ(OrderConsumer)
	// mq.MustStartPullConsumeBookProductMQ()
}
