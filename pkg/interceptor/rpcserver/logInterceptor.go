package rpcserver

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	zrpcErr "github.com/zeromicro/x/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// 对错误的日志信息记录，利用grpc拦截器实现

func LogInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any,
	err error) {
	resp, err = handler(ctx, req)
	if err == nil {
		return resp, nil
	}

	logx.WithContext(ctx).Errorf("【RPC SRV ERR】 %v", err)

	causeErr := errors.Cause(err)
	if e, ok := causeErr.(*zrpcErr.CodeMsg); ok {
		err = status.Error(codes.Code(e.Code), e.Msg)
	}

	return resp, err
}
