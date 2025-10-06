package logic

import (
	"context"
	"github.com/pkg/errors"

	"PaiPai/apps/user/models"
	"PaiPai/apps/user/rpc/internal/svc"
	"PaiPai/apps/user/rpc/user"
	"PaiPai/pkg/xerr"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

var ErrUserNotFound = xerr.New(xerr.SERVER_COMMON_ERROR, "该用户不存在")

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserInfoLogic) GetUserInfo(in *user.GetUserInfoReq) (*user.GetUserInfoResp, error) {
	// RPC层负责完整的业务验证和数据验证
	l.Logger.Info("[GetUserInfo] 获取用户信息请求", logx.Field("userId", in.Id))

	userEntiy, err := l.svcCtx.UsersModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == models.ErrNotFound {
			l.Logger.Error("[GetUserInfo] 用户不存在", logx.Field("userId", in.Id))
			return nil, ErrUserNotFound
		}
		l.Logger.Error("[GetUserInfo] 查询用户失败", logx.Field("userId", in.Id), logx.Field("err", err))
		return nil, errors.Wrapf(xerr.NewDBErr(), "find user by id err %v , req %v", err, in.Id)
	}

	l.Logger.Info("[GetUserInfo] 获取用户信息成功", logx.Field("userId", in.Id))
	var resp user.UserEntity
	copier.Copy(&resp, userEntiy) // copier.Copy复制结构体

	return &user.GetUserInfoResp{
		User: &resp,
	}, nil
}
