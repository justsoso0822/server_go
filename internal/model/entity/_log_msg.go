// =================================================================================
// 代码由 GoFrame CLI 工具生成并维护。请勿编辑。
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// LogMsg 是表 _log_msg 的 Go 结构体。
type LogMsg struct {
	Id   int         `json:"id"   orm:"id"   description:""`
	Uid  int         `json:"uid"  orm:"uid"  description:""`
	Msg  string      `json:"msg"  orm:"msg"  description:""`
	Time *gtime.Time `json:"time" orm:"time" description:""`
}
