package logic

import (
	"PaiPai/apps/user/rpc/internal/config"
	"PaiPai/apps/user/rpc/internal/svc"
	"github.com/zeromicro/go-zero/core/conf"
	"path/filepath"
)

// 单例测试

var svcCtx *svc.ServiceContext

func init() {
	var c config.Config
	conf.MustLoad(filepath.Join("../../etc/dev/user.yaml"), &c)
	svcCtx = svc.NewServiceContext(c)
}
