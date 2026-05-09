// =================================================================================
// 代码由 GoFrame CLI 工具生成并维护。请勿编辑。
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// LogMsg 是表 log_msg 的 Go 结构体，用于 Where/Data 等 DAO 操作。
type LogMsg struct {
	g.Meta `orm:"table:log_msg, do:true"`
	Id     any         //
	Uid    any         //
	Msg    any         //
	Time   *gtime.Time //
}
