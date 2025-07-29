package svc

import (
	model "PaiPai/apps/im/immodels"
	"PaiPai/apps/im/rpc/internal/config"
)

type ServiceContext struct {
	Config config.Config

	model.ChatLogModel
	model.ConversationsModel
	model.ConversationModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,

		ChatLogModel:       model.MustChatLogModel(c.Mongo.Url, c.Mongo.Db),
		ConversationsModel: model.MustConversationsModel(c.Mongo.Url, c.Mongo.Db),
		ConversationModel:  model.MustConversationModel(c.Mongo.Url, c.Mongo.Db),
	}
}
