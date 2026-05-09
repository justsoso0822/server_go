// =================================================================================
// 此文件由 GoFrame CLI 工具自动生成，可按需修改。
// =================================================================================

package dao

import (
	"server_go/internal/dao/internal"
)

// sysGmDao 是表 sys_gm 的数据访问对象。
// 可按需在此定义自定义方法以扩展功能。
type sysGmDao struct {
	*internal.SysGmDao
}

var (
	// SysGm 是表 sys_gm 操作的全局访问对象。
	SysGm = sysGmDao{internal.NewSysGmDao()}
)

// 在下方添加自定义方法和功能。
