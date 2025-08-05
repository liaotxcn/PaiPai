package job

import (
	"context"
	"github.com/pkg/errors"
	"time"
)

// RetryJetLagFunc 定义重试的时间策略
type RetryJetLagFunc func(ctx context.Context, retryCount int, lastTime time.Duration) time.Duration

func RetryJetLagAlways(ctx context.Context, retryCount int, lastTime time.Duration) time.Duration {
	return DefaultRetryJetLag
}

// IsRetryFunc 是否进行重试
type IsRetryFunc func(ctx context.Context, retryCount int, err error) bool

func RetryAlways(ctx context.Context, retryCount int, err error) bool {
	return true
}

func WithRetry(ctx context.Context, handler func(ctx context.Context) error, opts ...RetryOptions) error {
	opt := newOptions(opts...)

	// 判断程序本身是否设置了超时
	_, ok := ctx.Deadline()
	if !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, opt.timeout)
		defer cancel()
	}

	var (
		herr        error
		retryJetLag time.Duration
		ch          = make(chan error, 1)
	)

	for i := 0; i < opt.retryNums; i++ {
		go func() {
			ch <- handler(ctx)
		}()
		select {
		case herr = <-ch:
			if herr == nil {
				return nil
			}
			if !opt.isRetryFunc(ctx, i, herr) {
				return herr
			}
			retryJetLag = opt.retryJetLag(ctx, i, retryJetLag)
			time.Sleep(retryJetLag)
		case <-ctx.Done():
			return errors.New("time extend")
		}
	}
	return herr
}
