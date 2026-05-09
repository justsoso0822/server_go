// =================================================================================
// 此文件由 GoFrame CLI 工具自动生成，可按需修改。
// =================================================================================

package dao

import (
	"server_go/internal/dao/internal"
)

// LogMsgDao 是表 log_msg 的数据访问对象。
// 可按需在此定义自定义方法以扩展功能。
type LogMsgDao struct {
	*internal.LogMsgDao
}

var (
	// LogMsg 是表 log_msg 操作的全局访问对象。
	LogMsg = LogMsgDao{internal.NewLogMsgDao()}
)

// 在下方添加自定义方法和功能。
