package logic

import (
	"PaiPai/apps/im/rpc/imclient"
	"context"
	"github.com/jinzhu/copier"

	"PaiPai/apps/im/api/internal/svc"
	"PaiPai/apps/im/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetChatLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 根据用户获取聊天记录
func NewGetChatLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetChatLogLogic {
	return &GetChatLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetChatLogLogic) GetChatLog(req *types.ChatLogReq) (resp *types.ChatLogResp, err error) {
	// todo: add your logic here and delete this line
	data, err := l.svcCtx.GetChatLog(l.ctx, &imclient.GetChatLogReq{
		ConversationId: req.ConversationId,
		StartSendTime:  req.StartSendTime,
		EndSendTime:    req.EndSendTime,
		Count:          req.Count,
	})
	if err != nil {
		return nil, err
	}
	var res types.ChatLogResp
	copier.Copy(&res, &data)

	return &res, nil
}
