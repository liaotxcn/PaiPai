package logic

import (
	"PaiPai/apps/im/rpc/im"
	"PaiPai/apps/social/rpc/socialclient"
	"PaiPai/apps/user/rpc/user"
	"PaiPai/pkg/bitmap"
	"context"

	"PaiPai/apps/im/api/internal/svc"
	"PaiPai/apps/im/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetChatLogReadRecordsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取消息已读/未读记录
func NewGetChatLogReadRecordsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetChatLogReadRecordsLogic {
	return &GetChatLogReadRecordsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetChatLogReadRecordsLogic) GetChatLogReadRecords(req *types.GetChatLogReadRecordsReq) (resp *types.GetChatLogReadRecordsResp, err error) {
	// todo: add your logic here and delete this line

	chatLogs, err := l.svcCtx.ImClient.GetChatLog(l.ctx, &im.GetChatLogReq{
		MsgId: req.MsgId,
	})
	if err != nil || len(chatLogs.List) == 0 {
		return nil, err
	}

	var (
		chatLog = chatLogs.List[0]
		reads   = []string{chatLog.SendId}
		unReads []string
		ids     []string
	)

	// 分别设置已读未读
	switch status.ChatType(chatLog.ChatType) {
	case status.SingleChatType:
		if len(chatLog.ReadRecords) == 0 || chatLog.ReadRecords[0] == 0 {
			unReads = []string{chatLog.RecvId}
		} else {
			reads = append(reads, chatLog.RecvId)
		}
		ids = []string{chatLog.RecvId, chatLog.SendId}
	case status.GroupChatType:
		groupUsers, err := l.svcCtx.SocialClient.GroupUsers(l.ctx, &socialclient.GroupUsersReq{
			GroupId: chatLog.RecvId,
		})
		if err != nil {
			return nil, err
		}

		bitmaps := bitmap.Lood(chatLog.ReadRecords)
		for _, members := range groupUsers.List {
			ids = append(ids, members.UserId)

			if members.UserId == chatLog.SendId {
				continue
			}

			if bitmaps.IsSet(members.UserId) {
				reads = append(reads, members.UserId)
			} else {
				unReads = append(unReads, members.UserId)
			}
		}
	}

	userEntities, err := l.svcCtx.UserClient.FindUser(l.ctx, &user.FindUserReq{
		Ids: ids,
	})
	if err != nil {
		return nil, err
	}
	userEntitySet := make(map[string]*user.UserEntity, len(userEntities.Users))
	for i, entity := range userEntities.Users {
		userEntitySet[entity.Id] = userEntities.Users[i]
	}

	// 设置手机号码
	for i, read := range reads {
		if u := userEntitySet[read]; u != nil {
			reads[i] = u.Phone
		}
	}
	for i, unread := range unReads {
		if u := userEntitySet[unread]; u != nil {
			unReads[i] = u.Phone
		}
	}

	return &types.GetChatLogReadRecordsResp{
		Reads:   reads,
		UnReads: unReads,
	}, nil
}
