// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// LogTrace is the golang structure of table _log_trace for DAO operations like Where/Data.
type LogTrace struct {
	g.Meta `orm:"table:_log_trace, do:true"`
	Id     any         //
	Uid    any         //
	Type   any         //
	Num    any         //
	Before any         //
	After  any         //
	Reason any         //
	Time   *gtime.Time //
}
