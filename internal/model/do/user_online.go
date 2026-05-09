// =================================================================================
// 代码由 GoFrame CLI 工具生成并维护。请勿编辑。
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UserOnline 是表 user_online 的 Go 结构体，用于 Where/Data 等 DAO 操作。
type UserOnline struct {
	g.Meta   `orm:"table:user_online, do:true"`
	Uid      any         //
	Day      *gtime.Time //
	TmOnline any         //
	TmUpdate *gtime.Time //
}
