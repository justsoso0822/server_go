// =================================================================================
// 此文件由 GoFrame CLI 工具自动生成，可按需修改。
// =================================================================================

package dao

import (
	"server_go/internal/dao/internal"
)

// userResDao 是表 user_res 的数据访问对象。
// 可按需在此定义自定义方法以扩展功能。
type userResDao struct {
	*internal.UserResDao
}

var (
	// UserRes 是表 user_res 操作的全局访问对象。
	UserRes = userResDao{internal.NewUserResDao()}
)

// 在下方添加自定义方法和功能。
