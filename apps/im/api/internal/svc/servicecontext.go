package svc

import (
	"PaiPai/apps/im/api/internal/config"
	"PaiPai/apps/im/rpc/imclient"
	"PaiPai/apps/social/rpc/socialclient"
	"PaiPai/apps/user/rpc/userclient"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config  config.Config
	UserRpc userclient.User
	Social  socialclient.Social
	imclient.Im
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		UserRpc: userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		Social:  socialclient.NewSocial(zrpc.MustNewClient(c.UserRpc)),
		Im:      imclient.NewIm(zrpc.MustNewClient(c.ImRpc)),
	}
}
