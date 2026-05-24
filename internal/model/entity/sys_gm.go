// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SysGm is the golang structure for table sys_gm.
type SysGm struct {
	Uid  int         `json:"uid"  orm:"uid"  description:""`
	Tips string      `json:"tips" orm:"tips" description:""`
	Time *gtime.Time `json:"time" orm:"time" description:""`
}
