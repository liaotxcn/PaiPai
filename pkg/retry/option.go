package job

import "time"

type (
	RetryOptions func(opts *retryOptions)

	retryOptions struct {
		timeout     time.Duration
		retryNums   int
		isRetryFunc IsRetryFunc     //  是否重试
		retryJetLag RetryJetLagFunc // 重试间隔时间
	}
)

func newOptions(opts ...RetryOptions) *retryOptions {
	opt := &retryOptions{
		timeout:     DefaultRetryTimeout,
		retryNums:   DefaultRetryNums,
		isRetryFunc: RetryAlways,
		retryJetLag: RetryJetLagAlways,
	}
	for _, o := range opts {
		o(opt)
	}
	return opt
}

func WithRetryTime(timeout time.Duration) RetryOptions {
	return func(opts *retryOptions) {
		if timeout > 0 {
			opts.timeout = timeout
		}
	}
}
func WithRetryNum(nums int) RetryOptions {
	return func(opts *retryOptions) {
		if nums > 1 {
			opts.retryNums = nums
		}
	}
}
func WithRetryFunc(retryFunc IsRetryFunc) RetryOptions {
	return func(opts *retryOptions) {
		if retryFunc != nil {
			opts.isRetryFunc = retryFunc
		}
	}
}
func WithRetryJetLagFunc(retryFunc RetryJetLagFunc) RetryOptions {
	return func(opts *retryOptions) {
		if retryFunc != nil {
			opts.retryJetLag = retryFunc
		}
	}
}
