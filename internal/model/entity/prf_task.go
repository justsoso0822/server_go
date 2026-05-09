// =================================================================================
// 代码由 GoFrame CLI 工具生成并维护。请勿编辑。
// =================================================================================

package entity

// PrfTask 是表 prf_task 的 Go 结构体。
type PrfTask struct {
	Id        int    `json:"id"        orm:"id"         description:""`
	Ser       int    `json:"ser"       orm:"ser"        description:"对应工厂的ser"`
	Tid       string `json:"tid"       orm:"tid"        description:"任务道具id"`
	Npc       int    `json:"npc"       orm:"npc"        description:""`
	StartLoop int    `json:"startLoop" orm:"start_loop" description:"开启循环的第一个任务有效,每组里如果都是0,不开启循环"`
}
