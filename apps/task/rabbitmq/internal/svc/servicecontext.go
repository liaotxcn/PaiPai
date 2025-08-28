package svc

import (
	"PaiPai/apps/task/rabbitmq/internal/config"
	"PaiPai/apps/task/rabbitmq/internal/rabbit"
	"fmt"
)

type ServiceContext struct {
	Config   config.Config
	RabbitMQ *rabbit.RabbitMQ
	Producer *rabbit.Producer
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 创建RabbitMQ连接
	rmq, err := rabbit.NewRabbitMQ(&c.RabbitMQ)
	if err != nil {
		panic(fmt.Errorf("failed to create RabbitMQ connection: %w", err))
	}

	// 创建生产者
	producer := rabbit.NewProducer(rmq, "xchange")

	// 声明交换器
	if err := producer.DeclareExchange("direct"); err != nil {
		panic(fmt.Errorf("failed to declare exchange: %w", err))
	}

	return &ServiceContext{
		Config:   c,
		RabbitMQ: rmq,
		Producer: producer,
	}
}
