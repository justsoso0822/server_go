// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// UserTask is the golang structure of table user_task for DAO operations like Where/Data.
type UserTask struct {
	g.Meta `orm:"table:user_task, do:true"`
	Uid    any //
	Taskid any //
	Addtm  any // 添加时间
	Done   any // 是否完成
	Donetm any // 完成时间
}
