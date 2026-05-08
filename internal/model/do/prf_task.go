// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// PrfTask is the golang structure of table prf_task for DAO operations like Where/Data.
type PrfTask struct {
	g.Meta    `orm:"table:prf_task, do:true"`
	Id        any //
	Ser       any // 对应工厂的ser
	Tid       any // 任务道具id
	Npc       any //
	StartLoop any // 开启循环的第一个任务有效,每组里如果都是0,不开启循环
}
