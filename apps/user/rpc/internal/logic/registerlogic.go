package logic

import (
	"context"
	"regexp"
	"strings"
	"time"

	"PaiPai/apps/user/models"
	"PaiPai/apps/user/rpc/internal/svc"
	"PaiPai/apps/user/rpc/user"
	"PaiPai/pkg/ctxdata"
	"PaiPai/pkg/encrypt"
	"PaiPai/pkg/utils"
	"PaiPai/pkg/wuid"
	"PaiPai/pkg/xerr"

	errors "github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

var (
	// 用户名正则表达式：3-20位字母、数字、下划线
	regexUsername = regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)
)

var (
	ErrPhoneIsRegister   = xerr.New(xerr.SERVER_COMMON_ERROR, "该手机号已注册")
	ErrUsernameIsRegister = xerr.New(xerr.SERVER_COMMON_ERROR, "该用户名已被使用")
	ErrEmailIsRegister    = xerr.New(xerr.SERVER_COMMON_ERROR, "该邮箱已被注册")
	ErrInvalidUsername    = xerr.New(xerr.REQUEST_PARAM_ERROR, "用户名格式不正确")
	ErrInvalidEmail       = xerr.New(xerr.REQUEST_PARAM_ERROR, "邮箱格式不正确")
	ErrPasswordTooShort   = xerr.New(xerr.REQUEST_PARAM_ERROR, "密码长度不能少于6位")
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *user.RegisterReq) (*user.RegisterResp, error) {
	// RPC层负责完整的业务验证和数据验证
	l.Logger.Info("[Register] 用户注册请求参数", logx.Field("username", in.Username), logx.Field("email", in.Email))

	// 1. 验证用户名格式
	if len(in.Username) < 3 || len(in.Username) > 20 || !regexUsername.MatchString(in.Username) {
		l.Logger.Error("[Register] 用户名格式验证失败", logx.Field("username", in.Username))
		return nil, ErrInvalidUsername
	}

	// 2. 验证密码长度
	if len(in.Password) < 6 {
		l.Logger.Error("[Register] 密码长度验证失败")
		return nil, ErrPasswordTooShort
	}

	// 3. 验证邮箱格式
	if in.Email != "" && !strings.Contains(in.Email, "@") {
		l.Logger.Error("[Register] 邮箱格式验证失败", logx.Field("email", in.Email))
		return nil, ErrInvalidEmail
	}

	// 4. 验证手机号是否已注册
	if in.Phone != "" {
		userEntity, err := l.svcCtx.UsersModel.FindByPhone(l.ctx, in.Phone)
		if err != nil && err != models.ErrNotFound {
			l.Logger.Error("[Register] 查询手机号失败", logx.Field("phone", in.Phone), logx.Field("err", err))
			return nil, errors.Wrapf(xerr.NewDBErr(), "find user by phone err %v", err)
		}

		if userEntity != nil {
			l.Logger.Error("[Register] 手机号已注册", logx.Field("phone", in.Phone))
			return nil, ErrPhoneIsRegister
		}
	}

	// 5. 验证用户名是否已注册
	userEntity, err := l.svcCtx.UsersModel.FindByUsername(l.ctx, in.Username)
	if err != nil && err != models.ErrNotFound {
		l.Logger.Error("[Register] 查询用户名失败", logx.Field("username", in.Username), logx.Field("err", err))
		return nil, errors.Wrapf(xerr.NewDBErr(), "find user by username err %v", err)
	}

	if userEntity != nil {
		l.Logger.Error("[Register] 用户名已注册", logx.Field("username", in.Username))
		return nil, ErrUsernameIsRegister
	}

	// 6. 验证邮箱是否已注册（如果提供了邮箱）
	if in.Email != "" {
		userEntity, err := l.svcCtx.UsersModel.FindByEmail(l.ctx, in.Email)
		if err != nil && err != models.ErrNotFound {
			l.Logger.Error("[Register] 查询邮箱失败", logx.Field("email", in.Email), logx.Field("err", err))
			return nil, errors.Wrapf(xerr.NewDBErr(), "find user by email err %v", err)
		}

		if userEntity != nil {
			l.Logger.Error("[Register] 邮箱已注册", logx.Field("email", in.Email))
			return nil, ErrEmailIsRegister
		}
	}

	// 定义用户数据 - 直接创建Users结构体实例
	newUser := &models.Users{
		Id:        wuid.GenUid(l.svcCtx.Config.Mysql.DataSource),
		Username:  in.Username,
		Email:     in.Email,
		Avatar:    in.Avatar,
		Nickname:  in.Nickname,
		Phone:     in.Phone,
		Status:    utils.ConvertToInt8(0),
		Sex:       utils.ConvertToInt8(in.Sex),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if len(in.Password) > 0 {
		genPassword, err := encrypt.GenPasswordHash([]byte(in.Password))
		if err != nil {
			return nil, err
		}
		newUser.Password = string(genPassword)
	}

	_, err = l.svcCtx.UsersModel.Insert(l.ctx, newUser)
	if err != nil {
		return nil, err
	}

	// 生成token
	now := time.Now().Unix()
	token, err := ctxdata.GetJwtToken(l.svcCtx.Config.Jwt.AccessSecret, now, l.svcCtx.Config.Jwt.AccessExpire,
		newUser.Id)
	if err != nil {
		return nil, err
	}

	return &user.RegisterResp{
		Token:  token,
		Expire: now + l.svcCtx.Config.Jwt.AccessExpire,
	}, nil
}
