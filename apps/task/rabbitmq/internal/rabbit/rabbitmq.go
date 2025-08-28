package rabbit

import (
	"PaiPai/apps/task/rabbitmq/internal/config"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"sync"
)

type RabbitMQ struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
	config  *config.RabbitMQConfig
	mu      sync.Mutex
}

func NewRabbitMQ(cfg *config.RabbitMQConfig) (*RabbitMQ, error) {
	rmq := &RabbitMQ{
		config: cfg,
	}

	if err := rmq.connect(); err != nil {
		return nil, err
	}

	return rmq, nil
}

func (r *RabbitMQ) connect() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	url := fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
		r.config.Username,
		r.config.Password,
		r.config.Host,
		r.config.Port,
		r.config.VHost,
	)

	conn, err := amqp091.Dial(url)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return fmt.Errorf("failed to open channel: %w", err)
	}

	r.conn = conn
	r.channel = channel

	return nil
}

func (r *RabbitMQ) Channel() *amqp091.Channel {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.channel
}

func (r *RabbitMQ) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}

func (r *RabbitMQ) EnsureConnection() error {
	if r.conn == nil || r.conn.IsClosed() {
		return r.connect()
	}
	return nil
}
