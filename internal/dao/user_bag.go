// =================================================================================
// 此文件由 GoFrame CLI 工具自动生成，可按需修改。
// =================================================================================

package dao

import (
	"server_go/internal/dao/internal"
)

// userBagDao 是表 user_bag 的数据访问对象。
// 可按需在此定义自定义方法以扩展功能。
type userBagDao struct {
	*internal.UserBagDao
}

var (
	// UserBag 是表 user_bag 操作的全局访问对象。
	UserBag = userBagDao{internal.NewUserBagDao()}
)

// 在下方添加自定义方法和功能。
