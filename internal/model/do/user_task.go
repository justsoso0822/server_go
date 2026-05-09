// =================================================================================
// 代码由 GoFrame CLI 工具生成并维护。请勿编辑。
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// UserTask 是表 user_task 的 Go 结构体，用于 Where/Data 等 DAO 操作。
type UserTask struct {
	g.Meta `orm:"table:user_task, do:true"`
	Uid    any //
	Taskid any //
	Addtm  any // 添加时间
	Done   any // 是否完成
	Donetm any // 完成时间
}
