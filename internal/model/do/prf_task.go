// =================================================================================
// 代码由 GoFrame CLI 工具生成并维护。请勿编辑。
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// PrfTask 是表 prf_task 的 Go 结构体，用于 Where/Data 等 DAO 操作。
type PrfTask struct {
	g.Meta    `orm:"table:prf_task, do:true"`
	Id        any //
	Ser       any // 对应工厂的ser
	Tid       any // 任务道具id
	Npc       any //
	StartLoop any // 开启循环的第一个任务有效,每组里如果都是0,不开启循环
}
