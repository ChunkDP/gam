package mq

import (
	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQ(url string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMQ{
		conn:    conn,
		channel: channel,
	}, nil
}

func (r *RabbitMQ) Close() {
	r.channel.Close()
	r.conn.Close()
}

func (r *RabbitMQ) Publish(queueName string, message []byte) error {
	_, err := r.channel.QueueDeclare(
		queueName, // 队列名称
		true,      // 持久化
		false,     // 自动删除
		false,     // 独占
		false,     // 不等待
		nil,       // 参数
	)
	if err != nil {
		return err
	}

	err = r.channel.Publish(
		"",        // 交换机
		queueName, // 路由键
		false,     // 强制
		false,     // 立即
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)
	return err
}

func (r *RabbitMQ) Consume(queueName string) (<-chan amqp.Delivery, error) {
	_, err := r.channel.QueueDeclare(
		queueName, // 队列名称
		true,      // 持久化
		false,     // 自动删除
		false,     // 独占
		false,     // 不等待
		nil,       // 参数
	)
	if err != nil {
		return nil, err
	}

	messages, err := r.channel.Consume(
		queueName, // 队列名称
		"",        // 消费者
		true,      // 自动确认
		false,     // 独占
		false,     // 不等待
		false,     // 参数
		nil,
	)
	return messages, err
}
