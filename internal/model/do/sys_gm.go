// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SysGm is the golang structure of table sys_gm for DAO operations like Where/Data.
type SysGm struct {
	g.Meta `orm:"table:sys_gm, do:true"`
	Uid    any         //
	Tips   any         //
	Time   *gtime.Time //
}
