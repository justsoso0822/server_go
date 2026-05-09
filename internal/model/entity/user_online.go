// =================================================================================
// 代码由 GoFrame CLI 工具生成并维护。请勿编辑。
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UserOnline 是表 user_online 的 Go 结构体。
type UserOnline struct {
	Uid      int         `json:"uid"      orm:"uid"       description:""`
	Day      *gtime.Time `json:"day"      orm:"day"       description:""`
	TmOnline int         `json:"tmOnline" orm:"tm_online" description:""`
	TmUpdate *gtime.Time `json:"tmUpdate" orm:"tm_update" description:""`
}
