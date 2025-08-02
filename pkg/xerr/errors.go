/**
 * @author: dn-jinmin/dn-jinmin
 * @doc:
 */

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
	return New(ServerCommonError, ErrMsg(ServerCommonError))
}
