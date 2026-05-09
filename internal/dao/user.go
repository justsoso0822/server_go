// =================================================================================
// 此文件由 GoFrame CLI 工具自动生成，可按需修改。
// =================================================================================

package dao

import (
	"server_go/internal/dao/internal"
)

// userDao 是表 user 的数据访问对象。
// 可按需在此定义自定义方法以扩展功能。
type userDao struct {
	*internal.UserDao
}

var (
	// User 是表 user 操作的全局访问对象。
	User = userDao{internal.NewUserDao()}
)

// 在下方添加自定义方法和功能。
