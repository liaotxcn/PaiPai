package rabbit

import "time"

// 错误处理 && 重连机制
func WithRetry(fn func() error, maxRetries int, delay time.Duration) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		if err = fn(); err == nil {
			return nil
		}
		time.Sleep(delay)
	}
	return err
}
