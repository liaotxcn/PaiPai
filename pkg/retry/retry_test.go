package job

import (
	"context"
	"testing"
)

func TestWithRetry(t *testing.T) {
	type args struct {
		ctx     context.Context
		handler func(ctx context.Context) error
		opts    []RetryOptions
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := WithRetry(tt.args.ctx, tt.args.handler, tt.args.opts...); (err != nil) != tt.wantErr {
				t.Errorf("WithRetry() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
