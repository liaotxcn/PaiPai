package logic

import (
	"PaiPai/apps/user/models"
	"PaiPai/pkg/xerr"
	"context"

	"github.com/pkg/errors"

	"PaiPai/apps/user/rpc/internal/svc"
	"PaiPai/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindUserLogic {
	return &FindUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindUserLogic) FindUser(in *user.FindUserReq) (*user.FindUserResp, error) {
	var users []models.Users

	if in.Phone != "" {
		// 根据手机号查找单个用户
		user, err := l.svcCtx.UsersModel.FindByPhone(l.ctx, in.Phone)
		if err != nil {
			if err != models.ErrNotFound {
				return nil, errors.Wrapf(err, "failed to find user by phone: %s", in.Phone)
			}
			// 用户不存在时返回空数组
			users = []models.Users{}
		} else {
			users = []models.Users{*user}
		}
	} else if len(in.Ids) > 0 {
		// 根据ID列表查找多个用户
		userList, err := l.svcCtx.UsersModel.ListByIds(l.ctx, in.Ids)
		if err != nil {
			l.Errorw("failed to find users by IDs", logx.Field("error", err), logx.Field("ids", in.Ids))
			return nil, errors.Wrapf(err, "failed to find users by IDs: %v", in.Ids)
		}
		// 转换为值类型切片
		users = make([]models.Users, 0, len(userList))
		for _, u := range userList {
			users = append(users, *u)
		}
	} else if in.Name != "" {
		// 根据用户名查找多个用户
		userList, err := l.svcCtx.UsersModel.ListByName(l.ctx, in.Name)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to find users by name: %s", in.Name)
		}
		// 转换为值类型切片
		users = make([]models.Users, 0, len(userList))
		for _, u := range userList {
			users = append(users, *u)
		}
	} else {
		return nil, errors.WithStack(xerr.ParamError)
	}

	// 转换为UserEntity响应
	userEntities := make([]*user.UserEntity, 0, len(users))
	for _, u := range users {
		// 安全地处理可能为nil的指针字段
		status := int32(0)
		if u.Status != nil {
			status = int32(*u.Status)
		}
		sex := int32(0)
		if u.Sex != nil {
			sex = int32(*u.Sex)
		}

		userEntities = append(userEntities, &user.UserEntity{
			Id:       u.Id,
			Avatar:   u.Avatar,
			Nickname: u.Nickname,
			Phone:    u.Phone,
			Status:   status,
			Sex:      sex,
		})
	}

	return &user.FindUserResp{User: userEntities}, nil

}
