package friend

import (
	"PaiPai/apps/social/api/internal/svc"
	"PaiPai/apps/social/api/internal/types"
	"PaiPai/apps/social/rpc/social"
	"PaiPai/pkg/ctxdata"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendsOnlineLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewFriendsOnlineLogic 好友在线情况
func NewFriendsOnlineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendsOnlineLogic {
	return &FriendsOnlineLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendsOnlineLogic) FriendsOnline(req *types.FriendsOnlineReq) (resp *types.FriendsOnlineResp, err error) {
	// todo: add your logic here and delete this line
	uid := ctxdata.GetUId(l.ctx)

	friendList, err := l.svcCtx.Social.FriendList(l.ctx, &social.FriendListReq{
		UserId: uid,
	})
	if err != nil || len(friendList.List) == 0 {
		return &types.FriendsOnlineResp{}, err
	}

	// 在缓存中查询在线用户
	Fids := make([]string, 0, len(friendList.List))
	for _, friend := range friendList.List {
		Fids = append(Fids, friend.FriendUid)
	}
	onlines, err := l.svcCtx.Redis.Hgetall(status.REDIS_ONLINE_USER)
	if err != nil {
		return nil, err
	}
	resOnLineList := make(map[string]bool, len(Fids))
	for _, fid := range Fids {
		if _, ok := onlines[fid]; ok {
			resOnLineList[fid] = true
		} else {
			resOnLineList[fid] = false
		}

	}
	resp.OnlineList = resOnLineList
	return
}
