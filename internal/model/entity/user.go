// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// User is the golang structure for table user.
type User struct {
	Uid      uint        `json:"uid"      orm:"uid"      description:""`
	Platform string      `json:"platform" orm:"platform" description:""`
	Openid   string      `json:"openid"   orm:"openid"   description:""`
	Sid      uint        `json:"sid"      orm:"sid"      description:""`
	Name     string      `json:"name"     orm:"name"     description:""`
	Head     string      `json:"head"     orm:"head"     description:""`
	Born     *gtime.Time `json:"born"     orm:"born"     description:""`
}
