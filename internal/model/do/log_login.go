// =================================================================================
// 代码由 GoFrame CLI 工具生成并维护。请勿编辑。
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// LogLogin 是表 log_login 的 Go 结构体，用于 Where/Data 等 DAO 操作。
type LogLogin struct {
	g.Meta   `orm:"table:log_login, do:true"`
	Id       any         //
	Uid      any         //
	Platform any         //
	Time     *gtime.Time //
}
