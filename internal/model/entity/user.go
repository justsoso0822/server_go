// =================================================================================
// 代码由 GoFrame CLI 工具生成并维护。请勿编辑。
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// User 是表 user 的 Go 结构体。
type User struct {
	Uid      uint        `json:"uid"      orm:"uid"      description:""`
	Platform string      `json:"platform" orm:"platform" description:""`
	Openid   string      `json:"openid"   orm:"openid"   description:""`
	Sid      uint        `json:"sid"      orm:"sid"      description:""`
	Name     string      `json:"name"     orm:"name"     description:""`
	Head     string      `json:"head"     orm:"head"     description:""`
	Born     *gtime.Time `json:"born"     orm:"born"     description:""`
}
