package logic

import (
	"context"

	"PaiPai/apps/eino_chat/rpc/eino"
	"PaiPai/apps/eino_chat/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type InvokeAgentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewInvokeAgentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InvokeAgentLogic {
	return &InvokeAgentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 调用已注册的 Agent
func (l *InvokeAgentLogic) InvokeAgent(in *eino.InvokeAgentReq) (*eino.InvokeAgentResp, error) {
	// todo: add your logic here and delete this line

	return &eino.InvokeAgentResp{}, nil
}
