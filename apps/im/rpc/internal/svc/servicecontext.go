package svc

import (
        model "PaiPai/apps/im/immodels"
        "PaiPai/apps/im/rpc/internal/config"
        "github.com/zeromicro/go-zero/core/logx"
)

type ServiceContext struct {
        Config config.Config

        ChatLogModel       model.ChatLogModel
        ConversationsModel model.ConversationsModel
        ConversationModel  model.ConversationModel
}

func NewServiceContext(c config.Config) *ServiceContext {
        ctx := &ServiceContext{
                Config: c,
        }

        // 只有在 MongoDB 配置有效时才初始化模型
        if c.Mongo.Url != "" && c.Mongo.Db != "" {
                ctx.ChatLogModel = model.MustChatLogModel(c.Mongo.Url, c.Mongo.Db)
                ctx.ConversationsModel = model.MustConversationsModel(c.Mongo.Url, c.Mongo.Db)
                ctx.ConversationModel = model.MustConversationModel(c.Mongo.Url, c.Mongo.Db)
        } else {
                logx.Info("MongoDB not configured, skipping model initialization")
        }

        return ctx
}
