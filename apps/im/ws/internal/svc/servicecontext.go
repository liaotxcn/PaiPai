package svc

import (
	"PaiPai/apps/im/models"
	"PaiPai/apps/im/ws/internal/config"
	"PaiPai/apps/task/mq/mqclient"
)

type ServiceContext struct {
	Config config.Config

	models.ChatLogModel
	mqclient.MsgChatTransferClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:                c,
		MsgChatTransferClient: mqclient.NewMsgChatTransferClient(c.MsgChatTransfer.Addrs, c.MsgChatTransfer.Topic),
		ChatLogModel:          models.MustChatLogModel(c.Mongo.Url, c.Mongo.Db),
	}
}
