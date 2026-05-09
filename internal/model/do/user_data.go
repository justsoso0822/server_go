// =================================================================================
// 代码由 GoFrame CLI 工具生成并维护。请勿编辑。
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// UserData 是表 user_data 的 Go 结构体，用于 Where/Data 等 DAO 操作。
type UserData struct {
	g.Meta `orm:"table:user_data, do:true"`
	Uid    any //
	Key    any //
	Value  any //
}
