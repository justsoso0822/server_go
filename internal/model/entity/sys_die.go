// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SysDie is the golang structure for table sys_die.
type SysDie struct {
	Uid  int         `json:"uid"  orm:"uid"  description:""`
	Tips string      `json:"tips" orm:"tips" description:"封掉用户登陆时候的错误提示"`
	Time *gtime.Time `json:"time" orm:"time" description:""`
}
