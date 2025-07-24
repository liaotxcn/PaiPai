package logic

import (
	"PaiPai/apps/im/ws/internal/svc"
	"PaiPai/apps/im/ws/websocket"
	"PaiPai/apps/im/ws/ws"
	"PaiPai/pkg/wuid"
	"context"
	"time"
)

type Conversation struct {
	ctx context.Context
	srv *websocket.Server
	svc *svc.ServiceContext
}

func NewConversation(ctx context.Context, srv *websocket.Server, svc *svc.ServiceContext) *Conversation {
	return &Conversation{
		ctx: ctx,
		srv: srv,
		svc: svc,
	}
}

func (l *Conversation) SingleChat(data *ws.Chat, userId string) error {
	if data.ConversationId == "" {
		data.ConversationId = wuid.CombineId(userId, data.RecvId)
	}

	time.Sleep(time.Minute)
	// 记录消息
	chatLog := immodels.ChatLog{
		ConversationId: data.ConversationId,
		SendId:         userId,
		RecvId:         data.RecvId,
		ChatType:       data.ChatType,
		MsgFrom:        0,
		MsgType:        data.MType,
		MsgContent:     data.Content,
		SendTime:       time.Now().UnixNano(),
	}
	err := l.svc.ChatLogModel.Insert(l.ctx, &chatLog)

	return err
}
