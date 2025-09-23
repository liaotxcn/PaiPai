package group

import (
	"PaiPai/apps/social/api/internal/svc"
	"PaiPai/apps/social/api/internal/types"
	"PaiPai/apps/social/rpc/social"
	constants "PaiPai/pkg/constant"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupUserOnlineLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 群在线用户
func NewGroupUserOnlineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupUserOnlineLogic {
	return &GroupUserOnlineLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupUserOnlineLogic) GroupUserOnline(req *types.GroupUserOnlineReq) (resp *types.GroupUserOnlineResp, err error) {
	// 初始化响应对象
	resp = &types.GroupUserOnlineResp{
		OnlineList: make(map[string]bool),
	}

	groupUsers, err := l.svcCtx.Social.GroupUsers(l.ctx, &social.GroupUsersReq{
		GroupId: req.GroupId,
	})
	if err != nil {
		l.Errorf("查询群成员失败, groupId: %s, err: %v", req.GroupId, err)
		return resp, nil
	}

	if len(groupUsers.List) == 0 {
		return resp, nil
	}

	// 在缓存中查询在线用户
	Uids := make([]string, 0, len(groupUsers.List))
	for _, user := range groupUsers.List {
		Uids = append(Uids, user.UserId)
	}
	onlines, err := l.svcCtx.Redis.Hgetall(constants.REDIS_ONLINE_USER)
	if err != nil {
		l.Errorf("查询在线用户失败, err: %v", err)
		return resp, nil
	}

	resOnLineList := make(map[string]bool, len(Uids))
	for _, fid := range Uids {
		if _, ok := onlines[fid]; ok {
			resOnLineList[fid] = true
		} else {
			resOnLineList[fid] = false
		}
	}

	resp.OnlineList = resOnLineList
	return resp, nil
}
