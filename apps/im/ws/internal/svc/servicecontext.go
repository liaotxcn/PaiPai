package svc

import (
	model "PaiPai/apps/im/immodels"
	"PaiPai/apps/im/ws/internal/config"
	"PaiPai/apps/task/mq/mqclient"
)

type ServiceContext struct {
	Config config.Config

	ChatLogModel model.ChatLogModel
	mqclient.MsgChatTransferClient
	mqclient.MsgReadTransferClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:                c,
		MsgChatTransferClient: mqclient.NewMsgChatTransferClient(c.MsgChatTransfer.Addrs, c.MsgChatTransfer.Topic, nil),
		MsgReadTransferClient: mqclient.NewMsgReadTransferClient(c.MsgReadTransfer.Addrs, c.MsgReadTransfer.Topic, nil),
		ChatLogModel:          model.MustChatLogModel(c.Mongo.Url, c.Mongo.Db),
	}
}
