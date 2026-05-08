// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// LogMsg is the golang structure of table _log_msg for DAO operations like Where/Data.
type LogMsg struct {
	g.Meta `orm:"table:_log_msg, do:true"`
	Id     any         //
	Uid    any         //
	Msg    any         //
	Time   *gtime.Time //
}
