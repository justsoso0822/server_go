// =================================================================================
// 代码由 GoFrame CLI 工具生成并维护。请勿编辑。
// =================================================================================

package entity

// UserTask 是表 user_task 的 Go 结构体。
type UserTask struct {
	Uid    int `json:"uid"    orm:"uid"    description:""`
	Taskid int `json:"taskid" orm:"taskid" description:""`
	Addtm  int `json:"addtm"  orm:"addtm"  description:"添加时间"`
	Done   int `json:"done"   orm:"done"   description:"是否完成"`
	Donetm int `json:"donetm" orm:"donetm" description:"完成时间"`
}
