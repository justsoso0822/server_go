// =================================================================================
// 此文件由 GoFrame CLI 工具自动生成，可按需修改。
// =================================================================================

package dao

import (
	"server_go/internal/dao/internal"
)

// LogTraceDao 是表 log_trace 的数据访问对象。
// 可按需在此定义自定义方法以扩展功能。
type LogTraceDao struct {
	*internal.LogTraceDao
}

var (
	// LogTrace 是表 log_trace 操作的全局访问对象。
	LogTrace = LogTraceDao{internal.NewLogTraceDao()}
)

// 在下方添加自定义方法和功能。
