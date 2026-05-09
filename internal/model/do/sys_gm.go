// =================================================================================
// 代码由 GoFrame CLI 工具生成并维护。请勿编辑。
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SysGm 是表 sys_gm 的 Go 结构体，用于 Where/Data 等 DAO 操作。
type SysGm struct {
	g.Meta `orm:"table:sys_gm, do:true"`
	Uid    any         //
	Tips   any         //
	Time   *gtime.Time //
}
