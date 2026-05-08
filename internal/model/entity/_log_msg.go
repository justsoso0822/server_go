// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// LogMsg is the golang structure for table _log_msg.
type LogMsg struct {
	Id   int         `json:"id"   orm:"id"   description:""`
	Uid  int         `json:"uid"  orm:"uid"  description:""`
	Msg  string      `json:"msg"  orm:"msg"  description:""`
	Time *gtime.Time `json:"time" orm:"time" description:""`
}
