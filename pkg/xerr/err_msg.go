// @author: dn-jinmin/dn-jinmin
// @doc: 错误消息映射表

package xerr

var codeText = map[int]string{
	SERVER_COMMON_ERROR: "服务器异常，稍后再尝试",
	REQUEST_PARAM_ERROR: "请求参数有误",
	DB_ERROR:            "数据库繁忙，稍后再尝试",
}

// GetErrMsg 从错误码映射表获取错误消息
func GetErrMsg(errcode int) string {
	if msg, ok := codeText[errcode]; ok {
		return msg
	}
	return codeText[SERVER_COMMON_ERROR]
}
