package logic

import (
	"PaiPai/apps/user/models"
	"PaiPai/pkg/xerr"
	"context"
	"fmt"
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
	// todo: add your logic here and delete this line

	var users = make([]models.Users, 1) // 第一个位置留给 phone、name 查询
	var userEntities []*user.UserEntity

	if in.Phone != "" {
		err, _ := l.svcCtx.UsersModel.FindByPhone(l.ctx, in.Phone)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to find api by phone: %s", in.Phone)
		}
	} else if len(in.Ids) > 0 {
		users = nil
		err, _ := l.svcCtx.UsersModel.ListByIds(l.ctx, in.Ids)
		if err != nil {
			fmt.Printf("\n\n\n %v \n\n\n", err)
			return nil, errors.Wrapf(err, "failed to find users by IDs: %v", in.Ids)
		}
	} else if in.Name != "" {
		err, _ := l.svcCtx.UsersModel.ListByName(l.ctx, in.Name)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to find users by name: %s", in.Name)
		}
	} else {
		return nil, errors.WithStack(xerr.ParamError)
	}

	userEntities = make([]*user.UserEntity, len(users))

	for index, u := range users {
		userEntities[index] = &user.UserEntity{
			Id:       u.Id,
			Avatar:   u.Avatar,
			Nickname: u.Nickname,
			Phone:    u.Phone,
			Status:   int32(*u.Status),
			Sex:      int32(*u.Sex),
		}
	}

	return &user.FindUserResp{User: userEntities}, nil

}
