package user

import (
	"PaiPai/apps/user/rpc/user"
	"context"
	"errors"

	"PaiPai/apps/user/api/internal/svc"
	"PaiPai/apps/user/api/internal/types"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户注册
func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	// 1. API层负责基本参数验证：密码一致性检查和必填字段验证
	if req.Password != req.ConfirmPassword {
		return nil, errors.New("两次输入的密码不一致")
	}

	// 2. 验证必填字段
	if req.Username == "" {
		return nil, errors.New("用户名不能为空")
	}
	if req.Password == "" {
		return nil, errors.New("密码不能为空")
	}

	// 3. 调用RPC服务进行注册（业务逻辑和数据验证由RPC层负责）
	registerResp, err := l.svcCtx.User.Register(l.ctx, &user.RegisterReq{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		Phone:    req.Phone,
		Nickname: req.Nickname,
		Avatar:   req.Avatar,
		Sex:      int32(req.Sex),
	})
	if err != nil {
		return nil, err
	}

	var res types.RegisterResp
	copier.Copy(&res, registerResp)

	return &res, nil
}
