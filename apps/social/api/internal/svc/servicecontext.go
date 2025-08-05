package svc

import (
	"PaiPai/apps/im/rpc/imclient"
	"PaiPai/apps/social/api/internal/config"
	"PaiPai/apps/social/rpc/socialclient"
	"PaiPai/apps/user/rpc/userclient"
	"PaiPai/pkg/interceptor"
	"PaiPai/pkg/middleware"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config                config.Config
	IdempotenceMiddleware rest.Middleware // 幂等中间件
	UserRpc               userclient.User
	Social                socialclient.Social
	imclient.Im

	*redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {

	return &ServiceContext{
		Config:  c,
		UserRpc: userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		Social: socialclient.NewSocial(zrpc.MustNewClient(c.UserRpc,
			// 重试 && 幂等
			//zrpc.WithDialOption(grpc.WithDefaultServiceConfig(retryPolicy)),
			zrpc.WithUnaryClientInterceptor(interceptor.DefaultIdempotentClient),
		)),

		Im:                    imclient.NewIm(zrpc.MustNewClient(c.ImRpc)),
		Redis:                 redis.MustNewRedis(c.Redisx),
		IdempotenceMiddleware: middleware.NewIdempotenceMiddleware().Handler,
	}
}
