package logic

import (
	"PaiPai/apps/im/rpc/im"
	"context"

	"PaiPai/apps/im/api/internal/svc"
	"PaiPai/apps/im/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetUpUserConversationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 建立会话
func NewSetUpUserConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetUpUserConversationLogic {
	return &SetUpUserConversationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetUpUserConversationLogic) SetUpUserConversation(req *types.SetUpUserConversationReq) (resp *types.SetUpUserConversationResp, err error) {
	// todo: add your logic here and delete this line
	_, err = l.svcCtx.SetUpUserConversation(l.ctx, &im.SetUpUserConversationReq{
		SendId:   req.SendId,
		RecvId:   req.RecvId,
		ChatType: req.ChatType,
	})

	return
}
