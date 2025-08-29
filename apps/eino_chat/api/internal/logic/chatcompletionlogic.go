package logic

import (
	"context"

	"PaiPai/apps/eino_chat/api/internal/svc"
	"PaiPai/apps/eino_chat/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChatCompletionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 聊天补全
func NewChatCompletionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatCompletionLogic {
	return &ChatCompletionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChatCompletionLogic) ChatCompletion(req *types.ChatCompletionReq) (resp *types.ChatCompletionResp, err error) {
	// todo: add your logic here and delete this line

	return
}
