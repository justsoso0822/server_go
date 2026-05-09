// =================================================================================
// 此文件由 GoFrame CLI 工具自动生成，可按需修改。
// =================================================================================

package dao

import (
	"server_go/internal/dao/internal"
)

// userBagTpDao 是表 user_bag_tp 的数据访问对象。
// 可按需在此定义自定义方法以扩展功能。
type userBagTpDao struct {
	*internal.UserBagTpDao
}

var (
	// UserBagTp 是表 user_bag_tp 操作的全局访问对象。
	UserBagTp = userBagTpDao{internal.NewUserBagTpDao()}
)

// 在下方添加自定义方法和功能。
