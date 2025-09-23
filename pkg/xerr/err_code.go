/**
 * @author: dn-jinmin/dn-jinmin
 * @doc:
 */

// 常量的定义信息中可依据自己的项目规范定义
// 如：可以定义常量的前三位是业务码、后三位是功能码

package xerr

// 错误码定义
const (
	// 基础错误码
	SERVER_COMMON_ERROR = 100001 // 服务器通用错误
	REQUEST_PARAM_ERROR = 100002 // 请求参数错误
	DB_ERROR            = 100003 // 数据库错误
	
	// 业务错误码
	USER_ERROR    = 200001 // 用户相关错误
	FRIEND_ERROR  = 200101 // 好友相关错误
	GROUP_ERROR   = 200201 // 群组相关错误
)
