package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"
	"productshop/common"
	"productshop/datamodels"
	"productshop/services"
	"sync"

	"github.com/streadway/amqp"
)

// url格式：amqp://账号:密码 @RabbitMQ 服务器地址:端口号/vHost
var MQURL = "amqp://" + common.RMQUser + ":" + common.RMQPawsd + "@" + common.RMQHost + ":" + common.RMQPort + "/" + common.RMQVHost

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	// 队列名称
	QueueName string
	// 交换机
	Exchange string
	// key
	Key string
	// 连接信息
	MqURL string
	// 锁
	sync.Mutex
}

// 创建结构体实例
func NewRabbitMQ(queueName string, exchange string, key string) *RabbitMQ {
	rabbitmq := &RabbitMQ{
		QueueName: queueName,
		Exchange:  exchange,
		Key:       key,
		MqURL:     MQURL,
	}
	var err error
	rabbitmq.conn, err = amqp.Dial(rabbitmq.MqURL)
	rabbitmq.failOnErr(err, "创建连接错误")

	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "获取 channel 失败")

	return rabbitmq
}

// 断开 channel 和 connection 连接
func (r *RabbitMQ) Destroy() {
	r.channel.Close()
	r.conn.Close()
}

// 错误处理函数
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Printf("%s:%s", message, err)
		panic(fmt.Sprintf("%s:%s", message, err))
	}
}

// 创建 Simple 模式的生产者实例
func NewRabbitMQSimple(queueName string) *RabbitMQ {
	return NewRabbitMQ(queueName, "", "")
}

// Simple 生产者
func (r *RabbitMQ) PublishSimple(message string) error {
	r.Lock()
	defer r.Unlock()
	// 1.申请队列 如果队列不存在则会自动创建，如果存在则跳过创建
	_, err := r.channel.QueueDeclare(
		// 队列名称
		r.QueueName,
		// 是否持久化
		false,
		// 是否自动删除
		false,
		// 是否具有排他性
		false,
		// 是否阻塞
		false,
		// 额外属性
		nil,
	)
	if err != nil {
		return err
	}

	// 2.发送消息到队列中
	r.channel.Publish(
		// 交换机
		r.Exchange,
		// 队列名称
		r.QueueName,
		// /如果为true根据exhange会条件的队列那么会把发送的消息返回给发送者
		false,
		// 如果为 true，当 exchange 发送消息队列后发现队列上没有绑定消费者，则会把消息还给发送者
		false,
		// 发送的信息
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)

	return nil
}

// 消费消息
func (r *RabbitMQ) ConsumeSimple(orderService services.IOrderService, productService services.IProductService) {
	r.Lock()
	defer r.Unlock()
	// 1.申请队列 如果队列不存在则会自动创建，如果存在则跳过创建
	_, err := r.channel.QueueDeclare(
		// 队列名称
		r.QueueName,
		// 是否持久化
		false,
		// 是否自动删除
		false,
		// 是否具有排他性
		false,
		// 是否阻塞
		false,
		// 额外属性
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}

	// 2.接收消息
	msgs, err := r.channel.Consume(
		// 队列名称
		r.QueueName,
		// 用来区分多个消费者
		"",
		// 是否自动应答
		false,
		// 是否具有排他性
		false,
		// 如果设置为 true，表示不能将同一个 connection 中发送的消息传递给这个 connection 中的消费
		false,
		// 队列消费是否阻塞
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}

	forever := make(chan bool)
	// 启用协程处理消息
	r.channel.Qos(1, 0, false)

	go func() {
		for d := range msgs {
			log.Printf("接收到消息 Received a message: %s", d.Body)
			message := &datamodels.Message{}
			err := json.Unmarshal([]byte(d.Body), message)
			if err != nil {
				fmt.Println(err)
			}
			// 插入订单
			_, err = orderService.InsertOrderByMessage(message)
			if err != nil {
				fmt.Println(err)
			}
			// 扣除商品数量
			err = productService.SubNumber(message.ProductID, message.ProductNum)
			if err != nil {
				fmt.Println(err)
			}
			// true表示确认所有未确认的消息，false表示确认当前消息
			d.Ack(false)
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
