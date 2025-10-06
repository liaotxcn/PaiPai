package logic

import (
	"PaiPai/apps/user/models"
	"PaiPai/pkg/ctxdata"
	"PaiPai/pkg/encrypt"
	"PaiPai/pkg/xerr"
	"context"
	"github.com/pkg/errors"
	"time"

	"PaiPai/apps/user/rpc/internal/svc"
	"PaiPai/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrPhoneNotRegister = xerr.New(xerr.SERVER_COMMON_ERROR, "手机号尚未注册")
	ErrUserPwdError     = xerr.New(xerr.SERVER_COMMON_ERROR, "密码不正确")
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *user.LoginReq) (*user.LoginResp, error) {
	// RPC层负责完整的业务验证和数据验证
	l.Logger.Info("[Login] 用户登录请求参数", logx.Field("username", in.Username))

	// 验证用户是否注册(根据用户名)
	userEntity, err := l.svcCtx.UsersModel.FindByUsername(l.ctx, in.Username)
	if err != nil {
		if err == models.ErrNotFound {
			l.Logger.Error("[Login] 用户不存在", logx.Field("username", in.Username))
			return nil, ErrUserNotFound
		}
		l.Logger.Error("[Login] 查询用户失败", logx.Field("username", in.Username), logx.Field("err", err))
		return nil, errors.Wrapf(xerr.NewDBErr(), "find user by username err %v , req %v", err, in.Username)
	}

	// 密码验证
	if !encrypt.ValidatePassword(in.Password, userEntity.Password) {
		l.Logger.Error("[Login] 密码验证失败", logx.Field("username", in.Username))
		return nil, ErrUserPwdError
	}

	l.Logger.Info("[Login] 用户登录成功", logx.Field("username", in.Username), logx.Field("userId", userEntity.Id))

	// 生成token
	now := time.Now().Unix()
	token, err := ctxdata.GetJwtToken(l.svcCtx.Config.Jwt.AccessSecret, now, l.svcCtx.Config.Jwt.AccessExpire,
		userEntity.Id)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "ctxdata get jwt token err %v", err)
	}

	return &user.LoginResp{
		Id:     userEntity.Id,
		Token:  token,
		Expire: now + l.svcCtx.Config.Jwt.AccessExpire,
	}, nil
}
