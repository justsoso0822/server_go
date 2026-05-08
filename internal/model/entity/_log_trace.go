// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// LogTrace is the golang structure for table _log_trace.
type LogTrace struct {
	Id     int64       `json:"id"     orm:"id"     description:""`
	Uid    int         `json:"uid"    orm:"uid"    description:""`
	Type   string      `json:"type"   orm:"type"   description:""`
	Num    int         `json:"num"    orm:"num"    description:""`
	Before int         `json:"before" orm:"before" description:""`
	After  int         `json:"after"  orm:"after"  description:""`
	Reason string      `json:"reason" orm:"reason" description:""`
	Time   *gtime.Time `json:"time"   orm:"time"   description:""`
}
