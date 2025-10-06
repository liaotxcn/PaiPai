package models

import (
	"context"
	"fmt"
	
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UsersModel = (*customUsersModel)(nil)

// 缓存键前缀
const (
	cacheUsersUsernamePrefix = "cache:users:username:"
	cacheUsersEmailPrefix    = "cache:users:email:"
)

type (
	// UsersModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUsersModel.
	UsersModel interface {
		usersModel
		// 根据用户名查询用户
		FindByUsername(ctx context.Context, username string) (*Users, error)
		// 根据邮箱查询用户
		FindByEmail(ctx context.Context, email string) (*Users, error)
	}

	customUsersModel struct {
		*defaultUsersModel
	}
)

// NewUsersModel returns a model for the database table.
func NewUsersModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) UsersModel {
	return &customUsersModel{
		defaultUsersModel: newUsersModel(conn, c, opts...),
	}
}

// FindByUsername 根据用户名查询用户
func (m *customUsersModel) FindByUsername(ctx context.Context, username string) (*Users, error) {
	usersUsernameKey := fmt.Sprintf("%s%v", cacheUsersUsernamePrefix, username)
	var resp Users
	err := m.QueryRowCtx(ctx, &resp, usersUsernameKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `username` = ? limit 1", usersRows, m.defaultUsersModel.table)
		return conn.QueryRowCtx(ctx, v, query, username)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// FindByEmail 根据邮箱查询用户
func (m *customUsersModel) FindByEmail(ctx context.Context, email string) (*Users, error) {
	usersEmailKey := fmt.Sprintf("%s%v", cacheUsersEmailPrefix, email)
	var resp Users
	err := m.QueryRowCtx(ctx, &resp, usersEmailKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `email` = ? limit 1", usersRows, m.defaultUsersModel.table)
		return conn.QueryRowCtx(ctx, v, query, email)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
