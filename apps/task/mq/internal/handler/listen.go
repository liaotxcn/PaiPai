package handler

import (
	"PaiPai/apps/task/mq/internal/handler/msgTransfer"
	"PaiPai/apps/task/mq/internal/svc"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/service"
)

type Listen struct {
	svc *svc.ServiceContext
}

func NewListen(svc *svc.ServiceContext) *Listen {
	return &Listen{svc: svc}
}

func (l *Listen) Services() []service.Service {
	// 创建消息处理器
	msgReadHandler := msgTransfer.NewMsgReadTransfer(l.svc)
	msgChatHandler := msgTransfer.NewMsgChatTransfer(l.svc)

	return []service.Service{
		kq.MustNewQueue(l.svc.Config.MsgReadTransfer, kq.WithHandle(msgReadHandler.Consume)),

		// todo: 此处可以加载多个消费者
		kq.MustNewQueue(l.svc.Config.MsgChatTransfer, kq.WithHandle(msgChatHandler.Consume)),
	}
}
