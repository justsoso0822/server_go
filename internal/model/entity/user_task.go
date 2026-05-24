// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// UserTask is the golang structure for table user_task.
type UserTask struct {
	Uid    int `json:"uid"    orm:"uid"    description:""`
	Taskid int `json:"taskid" orm:"taskid" description:""`
	Addtm  int `json:"addtm"  orm:"addtm"  description:"添加时间"`
	Done   int `json:"done"   orm:"done"   description:"是否完成"`
	Donetm int `json:"donetm" orm:"donetm" description:"完成时间"`
}
