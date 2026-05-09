// =================================================================================
// 代码由 GoFrame CLI 工具生成并维护。请勿编辑。
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// UserBagTp 是表 user_bag_tp 的 Go 结构体，用于 Where/Data 等 DAO 操作。
type UserBagTp struct {
	g.Meta  `orm:"table:user_bag_tp, do:true"`
	Id      any //
	Uid     any //
	Chapter any // 副本id - 0 主线
	Time    any // 更新时间
	Itemid  any // 物品id, 0=空
	Count   any // 道具数量
}
