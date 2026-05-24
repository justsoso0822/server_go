// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UserOnline is the golang structure of table user_online for DAO operations like Where/Data.
type UserOnline struct {
	g.Meta   `orm:"table:user_online, do:true"`
	Uid      any         //
	Day      *gtime.Time //
	TmOnline any         //
	TmUpdate *gtime.Time //
}
