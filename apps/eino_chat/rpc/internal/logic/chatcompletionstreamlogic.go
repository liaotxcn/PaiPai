package logic

import (
	"context"

	"PaiPai/apps/eino_chat/rpc/eino"
	"PaiPai/apps/eino_chat/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChatCompletionStreamLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChatCompletionStreamLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatCompletionStreamLogic {
	return &ChatCompletionStreamLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 流式对话
func (l *ChatCompletionStreamLogic) ChatCompletionStream(in *eino.ChatCompletionStreamReq, stream eino.Eino_ChatCompletionStreamServer) error {
	// todo: add your logic here and delete this line

	return nil
}
