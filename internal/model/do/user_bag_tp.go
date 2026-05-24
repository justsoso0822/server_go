// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// UserBagTp is the golang structure of table user_bag_tp for DAO operations like Where/Data.
type UserBagTp struct {
	g.Meta  `orm:"table:user_bag_tp, do:true"`
	Id      any //
	Uid     any //
	Chapter any // 副本id - 0 主线
	Time    any // 更新时间
	Itemid  any // 物品id, 0=空
	Count   any // 道具数量
}
