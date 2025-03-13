package rabbitmq

import (
	"log"

	"github.com/streadway/amqp"
)

// RabbitMQ 封装RabbitMQ连接和操作
type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// NewRabbitMQ 创建RabbitMQ实例
func NewRabbitMQ(url string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &RabbitMQ{
		conn:    conn,
		channel: ch,
	}, nil
}

// Close 关闭连接
func (r *RabbitMQ) Close() {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
}

// DeclareQueue 声明队列
func (r *RabbitMQ) DeclareQueue(name string) (amqp.Queue, error) {
	return r.channel.QueueDeclare(
		name,  // 队列名
		true,  // 持久化
		false, // 自动删除
		false, // 排他性
		false, // 不等待
		nil,   // 参数
	)
}

// PublishMessage 发布消息
func (r *RabbitMQ) PublishMessage(queueName string, body []byte) error {
	_, err := r.DeclareQueue(queueName)
	if err != nil {
		return err
	}

	return r.channel.Publish(
		"",        // 交换机
		queueName, // 路由键
		false,     // 强制
		false,     // 立即
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

// ConsumeMessages 消费消息
func (r *RabbitMQ) ConsumeMessages(queueName string, handler func([]byte) error) error {
	q, err := r.DeclareQueue(queueName)
	if err != nil {
		return err
	}

	msgs, err := r.channel.Consume(
		q.Name, // 队列
		"",     // 消费者
		true,   // 自动确认
		false,  // 排他性
		false,  // 不等待
		false,  // 参数
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			if err := handler(d.Body); err != nil {
				log.Printf("Error handling message: %v", err)
			}
		}
	}()

	return nil
}
