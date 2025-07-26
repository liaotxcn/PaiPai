package svc

import (
	model "PaiPai/apps/im/immodels"
	"PaiPai/apps/im/ws/websocket"
	"PaiPai/apps/task/mq/internal/config"
	constants "PaiPai/pkg/constant"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"net/http"
)

type ServiceContext struct {
	config.Config

	WsClient websocket.Client
	*redis.Redis

	model.ChatLogModel
	model.ConversationModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	svc := &ServiceContext{
		Config:            c,
		Redis:             redis.MustNewRedis(c.Redisx),
		ChatLogModel:      model.MustChatLogModel(c.Mongo.Url, c.Mongo.Db),
		ConversationModel: model.MustConversationModel(c.Mongo.Url, c.Mongo.Db),
	}

	token, err := svc.GetSystemToken()
	if err != nil {
		panic(err)
	}

	header := http.Header{}
	header.Set("Authorization", token)
	svc.WsClient = websocket.NewClient(c.Ws.Host, websocket.WithClientHeader(header))
	return svc
}

func (svc *ServiceContext) GetSystemToken() (string, error) {
	return svc.Redis.Get(constants.REDIS_SYSTEM_ROOT_TOKEN)
}
