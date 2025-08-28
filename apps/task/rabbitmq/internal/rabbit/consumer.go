package rabbit

import (
	"context"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

type MessageHandler func(ctx context.Context, message []byte) error

type Consumer struct {
	rabbitMQ     *RabbitMQ
	queue        string
	exchange     string
	routingKey   string
	handler      MessageHandler
	RetryHandler func(ctx context.Context, delivery amqp091.Delivery, handler func(ctx context.Context, msg amqp091.Delivery) error) error
}

func NewConsumer(rabbitMQ *RabbitMQ, queue, exchange, routingKey string, handler MessageHandler) *Consumer {
	return &Consumer{
		rabbitMQ:   rabbitMQ,
		queue:      queue,
		exchange:   exchange,
		routingKey: routingKey,
		handler:    handler,
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	if err := c.rabbitMQ.EnsureConnection(); err != nil {
		return err
	}

	// 声明队列
	_, err := c.rabbitMQ.Channel().QueueDeclare(
		c.queue, // name
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	// 绑定队列到交换器
	err = c.rabbitMQ.Channel().QueueBind(
		c.queue,      // queue name
		c.routingKey, // routing key
		c.exchange,   // exchange
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to bind queue: %w", err)
	}

	// 消费消息
	msgs, err := c.rabbitMQ.Channel().Consume(
		c.queue, // queue
		"",      // consumer
		false,   // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	go c.consumeMessages(ctx, msgs)
	log.Printf("Consumer started for queue: %s", c.queue)
	return nil
}

func (c *Consumer) consumeMessages(ctx context.Context, msgs <-chan amqp091.Delivery) {
	for {
		select {
		case <-ctx.Done():
			log.Printf("Consumer stopped for queue: %s", c.queue)
			return
		case msg, ok := <-msgs:
			if !ok {
				log.Printf("Message channel closed for queue: %s", c.queue)
				return
			}

			c.processMessage(ctx, msg)
		}
	}
}

func (c *Consumer) processMessage(ctx context.Context, msg amqp091.Delivery) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic in message processing: %v", r)
			msg.Nack(false, true) // 重新入队
		}
	}()

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// 处理消息
	if err := c.handler(ctx, msg.Body); err != nil {
		log.Printf("Failed to process message: %v", err)
		// 处理失败，可以选择重新入队或放入死信队列
		msg.Nack(false, true) // 重新入队
		return
	}

	// 确认消息
	msg.Ack(false)
}
