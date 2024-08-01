# README
消费端采用 RabbitMQ 的消息推送机制，Broker 会不断地推送消息给消费者。不需要客户端主动来拉。

推送当消息的个数会受到 **channel.Qos** 的限制

```go
r.channel.Qos(1, 0, false)
```
1. prefetchCount (1)：表示消费者在确认前可以接收的最大消息数。在这个例子中，设置为 1 表示消费者在确认前最多接收一条消息。
2. prefetchSize (0)：表示消费者在确认前可以接收的最大内容大小（以字节为单位）。设置为 0 表示不限制内容大小。
3. global (false)：表示 QoS 设置是否应用于整个通道。如果设置为 true，QoS 设置将应用于通道上的所有消费者。如果设置为 false，QoS 设置将仅应用于单个消费者。