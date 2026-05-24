// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// UserRes is the golang structure of table user_res for DAO operations like Where/Data.
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
