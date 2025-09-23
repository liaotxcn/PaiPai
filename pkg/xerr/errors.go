// @author: dn-jinmin/dn-jinmin
// @doc: 错误处理工具包，提供统一的错误创建和处理方法

package xerr

import "github.com/zeromicro/x/errors"

// 统一封装错误创建方法

func New(code int, msg string) error {
	return errors.New(code, msg)
}

func NewMsg(msg string) error {
	return errors.New(SERVER_COMMON_ERROR, msg)
}

func NewDBErr() error {
	return errors.New(DB_ERROR, ErrMsg(DB_ERROR))
}

func NewServerCommonErr() error {
	return New(SERVER_COMMON_ERROR, ErrMsg(SERVER_COMMON_ERROR))
}

// ErrMsg 根据错误码获取错误信息
func ErrMsg(code int) string {
	switch code {
	case SERVER_COMMON_ERROR:
		return "服务器内部错误"
	case REQUEST_PARAM_ERROR:
		return "请求参数错误"
	case DB_ERROR:
		return "数据库操作错误"
	default:
		return "未知错误"
	}
}
