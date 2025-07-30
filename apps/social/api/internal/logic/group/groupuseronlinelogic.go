package group

import (
	"context"

	"PaiPai/apps/social/api/internal/svc"
	"PaiPai/apps/social/api/internal/types"

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

func (l *GroupUserOnlineLogic) GroupUserOnline(req *types.GroupUserOnlineReq) (resp *types.GroupUserListResp, err error) {
	// todo: add your logic here and delete this line

	groupUsers, err := l.svcCtx.Social.GroupUsers(l.ctx, &socialclient.GroupUsersReq{
		GroupId: req.GroupId,
	})
	if err != nil || len(groupUsers.List) == 0 {
		return
	}

	// 在缓存中查询在线用户
	Uids := make([]string, 0, len(groupUsers.List))
	for _, user := range groupUsers.List {
		Uids = append(Uids, user.UserId)
	}
	onlines, err := l.svcCtx.Redis.Hgetall(status.REDIS_ONLINE_USER)
	if err != nil {
		return nil, err
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
	return
}
