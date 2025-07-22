package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf

	UserRpc zrpc.RpcClientConf // 需调用user-rpc服务

	JwtAuth struct {
		AccessSecret string
		AccessExpire int64
	}
}
