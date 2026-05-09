// =================================================================================
// 代码由 GoFrame CLI 工具生成并维护。请勿编辑。
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// LogLogin 是表 log_login 的 Go 结构体。
type LogLogin struct {
	Id       int         `json:"id"       orm:"id"       description:""`
	Uid      int         `json:"uid"      orm:"uid"      description:""`
	Platform string      `json:"platform" orm:"platform" description:""`
	Time     *gtime.Time `json:"time"     orm:"time"     description:""`
}
