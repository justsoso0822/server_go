// =================================================================================
// 此文件由 GoFrame CLI 工具自动生成，可按需修改。
// =================================================================================

package dao

import (
	"server_go/internal/dao/internal"
)

// memConfigDao 是表 mem_config 的数据访问对象。
// 可按需在此定义自定义方法以扩展功能。
type memConfigDao struct {
	*internal.MemConfigDao
}

var (
	// MemConfig 是表 mem_config 操作的全局访问对象。
	MemConfig = memConfigDao{internal.NewMemConfigDao()}
)

// 在下方添加自定义方法和功能。
