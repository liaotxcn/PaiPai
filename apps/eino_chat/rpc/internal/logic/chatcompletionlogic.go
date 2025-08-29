package logic

import (
	"context"

	"PaiPai/apps/eino_chat/rpc/eino"
	"PaiPai/apps/eino_chat/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChatCompletionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChatCompletionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatCompletionLogic {
	return &ChatCompletionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 单次对话
func (l *ChatCompletionLogic) ChatCompletion(in *eino.ChatCompletionReq) (*eino.ChatCompletionResp, error) {
	// todo: add your logic here and delete this line

	return &eino.ChatCompletionResp{}, nil
}
