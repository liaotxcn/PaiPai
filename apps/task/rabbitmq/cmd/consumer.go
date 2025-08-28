package cmd

import (
	"PaiPai/apps/task/rabbitmq/handler"
	"PaiPai/apps/task/rabbitmq/internal/config"
	"PaiPai/apps/task/rabbitmq/internal/rabbit"
	"PaiPai/apps/task/rabbitmq/internal/svc"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 加载配置
	var c config.Config

	// 创建服务上下文
	svcCtx := svc.NewServiceContext(c)

	// 定义消息处理器
	messageHandler := func(ctx context.Context, body []byte) error {
		var message map[string]interface{}
		if err := json.Unmarshal(body, &message); err != nil {
			return fmt.Errorf("failed to unmarshal message: %w", err)
		}

		log.Printf("Received message: %+v", message)

		// 根据消息类型处理
		switch message["event"] {
		case "user_created":
			return handler.HandleCreated(ctx, message)
		default:
			log.Printf("Unknown event type: %s", message["event"])
		}

		return nil
	}

	// 创建消费者
	consumer := rabbit.NewConsumer(
		svcCtx.RabbitMQ,
		"user-queue",   // 队列名
		"exchange",     // 交换器
		"user.created", // 路由键
		messageHandler,
	)

	// 启动消费者
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := consumer.Start(ctx); err != nil {
		log.Fatalf("Failed to start consumer: %v", err)
	}

	// 等待中断信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down consumer...")
	cancel()
	svcCtx.RabbitMQ.Close()
}
