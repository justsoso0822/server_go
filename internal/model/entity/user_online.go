// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UserOnline is the golang structure for table user_online.
type UserOnline struct {
	Uid      int         `json:"uid"      orm:"uid"       description:""`
	Day      *gtime.Time `json:"day"      orm:"day"       description:""`
	TmOnline int         `json:"tmOnline" orm:"tm_online" description:""`
	TmUpdate *gtime.Time `json:"tmUpdate" orm:"tm_update" description:""`
}
