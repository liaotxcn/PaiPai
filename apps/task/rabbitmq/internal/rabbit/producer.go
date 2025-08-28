package rabbit

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

type Producer struct {
	rabbitMQ *RabbitMQ
	exchange string
}

func NewProducer(rabbitMQ *RabbitMQ, exchange string) *Producer {
	return &Producer{
		rabbitMQ: rabbitMQ,
		exchange: exchange,
	}
}

func (p *Producer) Publish(ctx context.Context, routingKey string, message interface{}) error {
	if err := p.rabbitMQ.EnsureConnection(); err != nil {
		return fmt.Errorf("failed to ensure connection: %w", err)
	}

	// 序列化消息
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// 发布消息
	err = p.rabbitMQ.Channel().PublishWithContext(
		ctx,
		p.exchange, // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp091.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp091.Persistent, // 持久化消息
			Timestamp:    time.Now(),
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Printf("Message published to exchange %s with routing key %s", p.exchange, routingKey)
	return nil
}

// 声明交换器
func (p *Producer) DeclareExchange(exchangeType string) error {
	if err := p.rabbitMQ.EnsureConnection(); err != nil {
		return err
	}

	return p.rabbitMQ.Channel().ExchangeDeclare(
		p.exchange,   // name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
}
