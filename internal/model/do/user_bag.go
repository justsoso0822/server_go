// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// UserBag is the golang structure of table user_bag for DAO operations like Where/Data.
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
