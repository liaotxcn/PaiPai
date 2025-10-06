package user

import (
	"PaiPai/apps/user/rpc/user"
	constants "PaiPai/pkg/constant"
	"context"
	"errors"
	"github.com/jinzhu/copier"

	"PaiPai/apps/user/api/internal/svc"
	"PaiPai/apps/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户登入
func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	// API层负责基本参数验证
	if req.Username == "" {
		return nil, errors.New("用户名不能为空")
	}
	if req.Password == "" {
		return nil, errors.New("密码不能为空")
	}

	// 调用RPC服务进行登录（登录业务逻辑和数据验证由RPC层负责）
	loginResp, err := l.svcCtx.User.Login(l.ctx, &user.LoginReq{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	var res types.LoginResp
	copier.Copy(&res, loginResp)

	// 记录用户在线状态
	l.svcCtx.Redis.HsetCtx(l.ctx, constants.REDIS_ONLINE_USER, loginResp.Id, "1")

	return &res, nil
}
