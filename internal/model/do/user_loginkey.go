// =================================================================================
// 代码由 GoFrame CLI 工具生成并维护。请勿编辑。
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// UserLoginkey 是表 user_loginkey 的 Go 结构体，用于 Where/Data 等 DAO 操作。
type UserLoginkey struct {
	g.Meta `orm:"table:user_loginkey, do:true"`
	Uid    any //
	Key    any //
	Ver    any // 客户端版本号
	Time   any // 登录时间
}
