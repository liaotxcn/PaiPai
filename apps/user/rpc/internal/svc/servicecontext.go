package svc

import (
	"PaiPai/apps/user/models"
	"PaiPai/apps/user/rpc/internal/config"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config     config.Config
	UserModels models.UsersModel
}

//func (c ServiceContext) SetRootToken() interface{} {
//
//}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.Mysql.DataSource)

	return &ServiceContext{
		Config:     c,
		UserModels: models.NewUsersModel(sqlConn, c.Cache),
	}
}
