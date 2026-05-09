// =================================================================================
// 代码由 GoFrame CLI 工具生成并维护。请勿编辑。
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// LogTrace 是表 log_trace 的 Go 结构体，用于 Where/Data 等 DAO 操作。
type LogTrace struct {
	g.Meta `orm:"table:log_trace, do:true"`
	Id     any         //
	Uid    any         //
	Type   any         //
	Num    any         //
	Before any         //
	After  any         //
	Reason any         //
	Time   *gtime.Time //
}
