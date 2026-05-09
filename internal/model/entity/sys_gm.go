// =================================================================================
// 代码由 GoFrame CLI 工具生成并维护。请勿编辑。
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SysGm 是表 sys_gm 的 Go 结构体。
type SysGm struct {
	Uid  int         `json:"uid"  orm:"uid"  description:""`
	Tips string      `json:"tips" orm:"tips" description:""`
	Time *gtime.Time `json:"time" orm:"time" description:""`
}
