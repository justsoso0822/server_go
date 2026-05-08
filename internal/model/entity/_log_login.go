// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// LogLogin is the golang structure for table _log_login.
type LogLogin struct {
	Id       int         `json:"id"       orm:"id"       description:""`
	Uid      int         `json:"uid"      orm:"uid"      description:""`
	Platform string      `json:"platform" orm:"platform" description:""`
	Time     *gtime.Time `json:"time"     orm:"time"     description:""`
}
