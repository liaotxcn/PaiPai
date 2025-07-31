package user

import (
	"PaiPai/apps/user/rpc/user"
	constants "PaiPai/pkg/constant"
	"context"
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
	// todo: add your logic here and delete this line

	loginResp, err := l.svcCtx.User.Login(l.ctx, &user.LoginReq{
		Phone:    req.Phone,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	var res types.LoginResp
	copier.Copy(&res, loginResp)

	// 处理登录业务
	l.svcCtx.Redis.HsetCtx(l.ctx, constants.REDIS_ONLINE_USER, loginResp.Id, "1")

	return &res, nil
}
