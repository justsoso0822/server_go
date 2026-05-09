// =================================================================================
// 代码由 GoFrame CLI 工具生成并维护。请勿编辑。
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// User 是表 user 的 Go 结构体，用于 Where/Data 等 DAO 操作。
type User struct {
	g.Meta   `orm:"table:user, do:true"`
	Uid      any         //
	Platform any         //
	Openid   any         //
	Sid      any         //
	Name     any         //
	Head     any         //
	Born     *gtime.Time //
}
