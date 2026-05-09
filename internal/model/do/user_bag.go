// =================================================================================
// 代码由 GoFrame CLI 工具生成并维护。请勿编辑。
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// UserBag 是表 user_bag 的 Go 结构体，用于 Where/Data 等 DAO 操作。
type UserBag struct {
	g.Meta  `orm:"table:user_bag, do:true"`
	Id      any //
	Uid     any //
	Chapter any // 副本id - 0-主格子
	Time    any // 操作更新时间
	Itemid  any // 物品id, 0=空
	Info    any // 格子信息
	Type    any // 格子类型,0-普通,1-道具购买
}
