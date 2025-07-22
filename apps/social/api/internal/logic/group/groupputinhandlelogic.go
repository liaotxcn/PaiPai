package group

import (
	"PaiPai/apps/social/rpc/socialclient"
	constants "PaiPai/pkg/constant"
	"PaiPai/pkg/ctxdata"
	"context"

	"PaiPai/apps/social/api/internal/svc"
	"PaiPai/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupPutInHandleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 申请进群处理
func NewGroupPutInHandleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutInHandleLogic {
	return &GroupPutInHandleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupPutInHandleLogic) GroupPutInHandle(req *types.GroupPutInHandleRep) (resp *types.GroupPutInHandleResp, err error) {
	// todo: add your logic here and delete this line

	_, err = l.svcCtx.Social.GroupPutInHandle(l.ctx, &socialclient.GroupPutInHandleReq{
		GroupReqId:   req.GroupReqId,
		GroupId:      req.GroupId,
		HandleUid:    ctxdata.GetUId(l.ctx),
		HandleResult: req.HandleResult,
	})

	if constants.HandlerResult(req.HandleResult) != constants.PassHandlerResult {
		return
	}

	// todo: 通过后的业务

	return
}
