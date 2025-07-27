package logic

import (
	model "PaiPai/apps/im/immodels"
	"PaiPai/pkg/constant"
	"PaiPai/pkg/xerr"
	"context"
	"github.com/pkg/errors"

	"PaiPai/apps/im/rpc/im"
	"PaiPai/apps/im/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateGroupConversationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateGroupConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateGroupConversationLogic {
	return &CreateGroupConversationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateGroupConversationLogic) CreateGroupConversation(in *im.CreateGroupConversationReq) (*im.CreateGroupConversationResp, error) {
	// todo: add your logic here and delete this line

	res := &im.CreateGroupConversationResp{}

	// 对群会话验证
	_, err := l.svcCtx.ConversationModel.FindOne(l.ctx, in.GroupId)
	if err != nil {
		return res, nil
	}
	if err != model.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewDBErr(), "ConversationModel.FindOne err %v, req %v", err, in.GroupId)
	}

	err = l.svcCtx.ConversationModel.Insert(l.ctx, &model.Conversation{
		ConversationId: in.GroupId,
		ChatType:       constants.GroupChatType,
	})

	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "ConversationModel.FindOne err %v, req %v", err)
	}

	_, err = NewSetUpUserConversationLogic(l.ctx, l.svcCtx).SetUpUserConversation(&im.SetUpUserConversationReq{
		SendId:   in.CreateId,
		RecvId:   in.GroupId,
		ChatType: int32(constants.GroupChatType),
	})

	return res, nil
}
