// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SysDie is the golang structure of table sys_die for DAO operations like Where/Data.
type SysDie struct {
	g.Meta `orm:"table:sys_die, do:true"`
	Uid    any         //
	Tips   any         // 封掉用户登陆时候的错误提示
	Time   *gtime.Time //
}
