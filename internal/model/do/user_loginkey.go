// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// UserLoginkey is the golang structure of table user_loginkey for DAO operations like Where/Data.
type UserLoginkey struct {
	g.Meta `orm:"table:user_loginkey, do:true"`
	Uid    any //
	Key    any //
	Ver    any // 客户端版本号
	Time   any // 登录时间
}
