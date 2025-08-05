package job

import "time"

const (
	DefaultRetryJetLag  = time.Second     // 默认重试超时时间
	DefaultRetryTimeout = time.Second * 2 // 任务超时时间
	DefaultRetryNums    = 3               // 默认重试次数
)
