// =================================================================================
// 代码由 GoFrame CLI 工具生成并维护。请勿编辑。
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// UserRes 是表 user_res 的 Go 结构体，用于 Where/Data 等 DAO 操作。
type UserRes struct {
	g.Meta   `orm:"table:user_res, do:true"`
	Uid      any //
	Gold     any //
	Diamond  any //
	Star     any //
	Tili     any //
	TiliTime any //
	Exp      any //
	Level    any //
	DayConf  any // 每日重置的数据
	DayTime  any // 上次重置时间
}
