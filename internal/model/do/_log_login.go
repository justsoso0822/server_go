// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// LogLogin is the golang structure of table _log_login for DAO operations like Where/Data.
type LogLogin struct {
	g.Meta   `orm:"table:_log_login, do:true"`
	Id       any         //
	Uid      any         //
	Platform any         //
	Time     *gtime.Time //
}
