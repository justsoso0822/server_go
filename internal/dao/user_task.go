// =================================================================================
// 此文件由 GoFrame CLI 工具自动生成，可按需修改。
// =================================================================================

package dao

import (
	"server_go/internal/dao/internal"
)

// userTaskDao 是表 user_task 的数据访问对象。
// 可按需在此定义自定义方法以扩展功能。
type userTaskDao struct {
	*internal.UserTaskDao
}

var (
	// UserTask 是表 user_task 操作的全局访问对象。
	UserTask = userTaskDao{internal.NewUserTaskDao()}
)

// 在下方添加自定义方法和功能。
