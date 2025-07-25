package config

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type Config struct {
	service.ServiceConf
	ListenOn string

	MsgChatTransfer kq.KqConf

	Redisx redis.RedisConf
	Mongo  struct {
		Url string
		Db  string
	}

	Ws struct {
		Host string
	}
}
