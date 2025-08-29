package logic

import (
	"context"

	"PaiPai/apps/eino_chat/api/internal/svc"
	"PaiPai/apps/eino_chat/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChatStreamLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 聊天流式输出
func NewChatStreamLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatStreamLogic {
	return &ChatStreamLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChatStreamLogic) ChatStream(req *types.ChatStreamReq) (resp *types.ChatStreamResp, err error) {
	// todo: add your logic here and delete this line

	return
}
