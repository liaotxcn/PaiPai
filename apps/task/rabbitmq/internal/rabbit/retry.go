package rabbit

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQ 重试机制
//	分级重试策略：不同次数不同处理方式
//  智能延迟：第一次立即重试，第二次30秒，后续指数退避
//  诊断信息：从第二次失败开始添加详细诊断信息
//  分级告警：根据失败次数触发不同级别的告警

// RetryConfig 重试配置
// 用于配置消息重试的各种参数，支持灵活的定制化配置
type RetryConfig struct {
	MaxRetries        int           // 最大重试次数，0表示不重试
	InitialDelay      time.Duration // 初始延迟时间，用于指数退避计算
	MaxDelay          time.Duration // 最大延迟时间，防止延迟时间过长
	BackoffFactor     float64       // 退避因子，每次重试延迟时间的增长倍数
	AlertThreshold    int           // 告警阈值，达到此重试次数触发轻度告警
	CriticalThreshold int           // 严重告警阈值，达到此重试次数触发严重告警
}

// DefaultRetryConfig 默认重试配置
// 返回一个合理的默认配置，适用于大多数场景
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries:        5,               // 最多重试5次
		InitialDelay:      1 * time.Second, // 初始延迟1秒
		MaxDelay:          5 * time.Minute, // 最大延迟5分钟
		BackoffFactor:     2.0,             // 每次延迟时间翻倍
		AlertThreshold:    2,               // 第2次失败触发告警
		CriticalThreshold: 3,               // 第3次失败触发严重告警
	}
}

// DeliveryWithRetry 带重试的消息投递器
// 封装了完整的重试逻辑，提供分级重试能力
type DeliveryWithRetry struct {
	config RetryConfig // 重试配置
}

// NewDeliveryWithRetry 创建新的重试投递器
// 参数: config - 重试配置，可以使用DefaultRetryConfig()获取默认配置
// 返回: *DeliveryWithRetry - 初始化好的重试投递器实例
func NewDeliveryWithRetry(config RetryConfig) *DeliveryWithRetry {
	return &DeliveryWithRetry{config: config}
}

// RetryDelivery 分级重试投递 - 核心重试逻辑
// 实现分级重试策略，根据重试次数采取不同的处理方式
//
// 参数:
//
//	ctx - 上下文，用于控制重试过程的取消
//	delivery - RabbitMQ消息投递对象，包含消息内容和元数据
//	handler - 消息处理函数，实际执行业务逻辑的函数
//
// 返回:
//
//	error - 如果所有重试都失败，返回最终错误；成功返回nil
func (d *DeliveryWithRetry) RetryDelivery(
	ctx context.Context,
	delivery amqp.Delivery,
	handler func(ctx context.Context, msg amqp.Delivery) error,
) error {
	var lastErr error
	// 从消息头中获取当前重试次数，支持从之前的重试中恢复
	retryCount := getRetryCount(delivery.Headers)

	// 重试循环：从当前重试次数+1开始，直到达到最大重试次数
	for attempt := retryCount + 1; attempt <= d.config.MaxRetries; attempt++ {
		// 执行实际的消息处理逻辑
		err := handler(ctx, delivery)
		if err == nil {
			log.Printf("Message processed successfully on attempt %d", attempt)
			return nil // 处理成功，立即返回
		}

		lastErr = err
		// 处理本次重试尝试：更新头信息、记录日志、触发告警等
		d.handleRetryAttempt(attempt, err, &delivery)

		// 计算下一次重试的延迟时间，采用智能延迟策略
		delay := d.calculateDelay(attempt)

		// 等待延迟时间，同时监听上下文取消信号
		select {
		case <-time.After(delay):
			continue // 延迟结束后继续下一次重试
		case <-ctx.Done():
			return ctx.Err() // 上下文被取消，立即返回
		}
	}

	// 所有重试尝试都失败，返回最终错误
	return fmt.Errorf("message failed after %d attempts: %w", d.config.MaxRetries, lastErr)
}

// handleRetryAttempt 处理单次重试尝试 - 分级处理核心
// 根据重试次数采取不同的处理策略：日志记录、告警、头信息更新等
//
// 参数:
//
//	attempt - 当前重试次数（从1开始）
//	err - 本次重试的错误信息
//	delivery - 消息投递对象，用于更新头信息
func (d *DeliveryWithRetry) handleRetryAttempt(attempt int, err error, delivery *amqp.Delivery) {
	// 更新消息头中的重试信息和诊断数据
	d.updateDeliveryHeaders(attempt, err, delivery)

	// 分级处理策略：不同重试次数采取不同措施
	switch attempt {
	case 1:
		// 第一次失败：通常是由于瞬时网络问题，立即重试
		// 记录日志但不告警，避免噪声
		log.Printf("First delivery attempt failed (will retry immediately): %v", err)
		// 立即重试，不需要额外延迟

	case 2:
		// 第二次失败：可能不是瞬时问题，需要关注
		// 延迟30秒后重试，触发轻度告警，添加诊断信息
		log.Printf("Second delivery attempt failed (will retry after 30s): %v", err)
		d.sendMildAlert(delivery, err) // 触发轻度告警

	case 3:
		// 第三次失败：问题比较严重，需要重点关注
		log.Printf("Third delivery attempt failed: %v", err)
		d.sendMediumAlert(delivery, err) // 触发中度告警

	default:
		// 第四次及以上：严重问题，可能需要人工干预
		log.Printf("Delivery attempt %d failed: %v", attempt, err)
		d.sendCriticalAlert(delivery, err, attempt) // 触发严重告警
	}
}

// calculateDelay 计算重试延迟 - 智能延迟策略
// 实现分级延迟：第一次立即，第二次固定30秒，后续指数退避
//
// 参数:
//
//	attempt - 当前重试次数
//
// 返回:
//
//	time.Duration - 计算出的延迟时间
func (d *DeliveryWithRetry) calculateDelay(attempt int) time.Duration {
	if attempt == 1 {
		return 0 // 第一次立即重试：处理瞬时网络问题
	}
	if attempt == 2 {
		return 30 * time.Second // 第二次固定30秒：给系统恢复时间
	}

	// 指数退避算法：延迟时间 = 初始延迟 * (退避因子 ^ (重试次数-2))
	// 避免重试过于频繁，同时保证不会无限延迟
	delay := time.Duration(float64(d.config.InitialDelay) *
		pow(d.config.BackoffFactor, float64(attempt-2)))

	// 限制最大延迟时间，防止延迟过长影响系统响应
	if delay > d.config.MaxDelay {
		return d.config.MaxDelay
	}
	return delay
}

// updateDeliveryHeaders 更新消息头中的诊断信息 - 诊断数据收集
// 在消息头中添加重试相关的元数据，用于问题排查和监控
//
// 参数:
//
//	attempt - 当前重试次数
//	err - 错误信息
//	delivery - 消息投递对象
func (d *DeliveryWithRetry) updateDeliveryHeaders(attempt int, err error, delivery *amqp.Delivery) {
	// 确保头信息map已初始化
	if delivery.Headers == nil {
		delivery.Headers = make(amqp.Table)
	}

	// 更新重试次数和时间戳
	delivery.Headers["x-retry-count"] = attempt
	delivery.Headers["x-last-retry-time"] = time.Now().Format(time.RFC3339)

	// 保存错误信息，但只保留最近3次错误避免头信息过大
	if attempt <= 3 {
		errorKey := fmt.Sprintf("x-error-attempt-%d", attempt)
		delivery.Headers[errorKey] = err.Error()
	}

	// 从第二次失败开始添加详细的诊断信息
	if attempt >= 2 {
		delivery.Headers["x-diagnostic-info"] = d.generateDiagnosticInfo(delivery, err)
	}
}

// generateDiagnosticInfo 生成诊断信息 - 问题排查数据
// 创建结构化的诊断信息，包含消息上下文和环境信息
//
// 参数:
//
//	delivery - 消息投递对象
//	err - 错误信息
//
// 返回:
//
//	string - JSON格式的诊断信息
func (d *DeliveryWithRetry) generateDiagnosticInfo(delivery *amqp.Delivery, err error) string {
	info := map[string]interface{}{
		"timestamp":    time.Now().Format(time.RFC3339Nano), // 精确时间戳
		"error":        err.Error(),                         // 错误详情
		"message_id":   delivery.MessageId,                  // 消息ID用于追踪
		"content_type": delivery.ContentType,                // 消息内容类型
		"routing_key":  delivery.RoutingKey,                 // 路由键
		"retry_count":  getRetryCount(delivery.Headers),     // 当前重试次数
		"host":         getHostInfo(),                       // 主机信息用于定位问题机器
	}

	jsonData, _ := json.Marshal(info)
	return string(jsonData)
}

// 告警方法 - 分级告警系统

// sendMildAlert 发送轻度告警
// 第二次失败时触发，用于提醒关注但不需要立即处理
func (d *DeliveryWithRetry) sendMildAlert(delivery *amqp.Delivery, err error) {
	// 实际项目中可以集成邮件、短信、钉钉、企业微信等告警渠道
	log.Printf("🚨 MILD ALERT: Message delivery failed twice. MessageID: %s, Error: %v",
		delivery.MessageId, err)
}

// sendMediumAlert 发送中度告警
// 第三次失败时触发，需要重点关注
func (d *DeliveryWithRetry) sendMediumAlert(delivery *amqp.Delivery, err error) {
	// 可以触发更高级别的通知，如电话告警
	log.Printf("🚨🚨 MEDIUM ALERT: Message delivery failed 3 times. MessageID: %s",
		delivery.MessageId)
}

// sendCriticalAlert 发送严重告警
// 第四次及以上失败时触发，可能需要立即人工干预
func (d *DeliveryWithRetry) sendCriticalAlert(delivery *amqp.Delivery, err error, attempt int) {
	// 触发最高级别的告警，确保相关人员立即知晓
	log.Printf("🚨🚨🚨 CRITICAL ALERT: Message delivery failed %d times. MessageID: %s, Error: %v",
		attempt, delivery.MessageId, err)
}

// 辅助函数

// getRetryCount 从消息头中获取重试次数
// 支持多种数据类型转换，确保兼容性
func getRetryCount(headers amqp.Table) int {
	if count, ok := headers["x-retry-count"].(int32); ok {
		return int(count)
	}
	if count, ok := headers["x-retry-count"].(int); ok {
		return count
	}
	return 0 // 默认返回0，表示第一次重试
}

// pow 简单的幂运算函数
// 用于指数退避计算
func pow(x, y float64) float64 {
	result := 1.0
	for i := 0; i < int(y); i++ {
		result *= x
	}
	return result
}

// getHostInfo 获取主机信息
// 实际项目中可以获取真实的主机名、IP、容器ID等信息
func getHostInfo() string {
	// 示例实现，实际中可以调用os.Hostname()等函数
	return "unknown-host"
}

// 使用示例
func ExampleUsage() {
	// 创建自定义重试配置
	config := RetryConfig{
		MaxRetries:        5,
		InitialDelay:      1 * time.Second,
		MaxDelay:          2 * time.Minute,
		BackoffFactor:     2.0,
		AlertThreshold:    2,
		CriticalThreshold: 3,
	}

	// 创建重试投递器实例
	retryDelivery := NewDeliveryWithRetry(config)

	// 在消费者中使用
	consumer := &Consumer{
		RetryHandler: retryDelivery.RetryDelivery, // 注入重试处理器
	}

	_ = consumer // 使用consumer
}
