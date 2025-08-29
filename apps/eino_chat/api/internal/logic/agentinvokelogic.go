package logic

import (
	"context"

	"PaiPai/apps/eino_chat/api/internal/svc"
	"PaiPai/apps/eino_chat/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AgentInvokeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 调用智能体
func NewAgentInvokeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AgentInvokeLogic {
	return &AgentInvokeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AgentInvokeLogic) AgentInvoke(req *types.AgentInvokeReq) (resp *types.AgentInvokeResp, err error) {
	// todo: add your logic here and delete this line

	return
}
