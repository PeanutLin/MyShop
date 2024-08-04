package main

import (
	"fmt"
	"productshop/db"
	"productshop/rabbitmq"
	"productshop/repositories"
	"productshop/services"
)

func main() {
	db, err := db.NewMysqlConn()
	if err != nil {
		fmt.Println(err)
	}

	// product 服务
	productRepository := repositories.NewProductManager(db)
	productService := services.NewProductService(productRepository)

	// order 服务
	orderRepository := repositories.NewOrderManagerRepository(db)
	orderService := services.NewOrderService(orderRepository)

	// rabitmq 简单模式
	rabbitmqConsumerSimple := rabbitmq.NewRabbitMQSimple("productshop")

	rabbitmqConsumerSimple.ConsumeSimple(orderService, productService)
}
