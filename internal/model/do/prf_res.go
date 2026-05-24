// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// PrfRes is the golang structure of table prf_res for DAO operations like Where/Data.
type PrfRes struct {
	g.Meta `orm:"table:prf_res, do:true"`
	Id     any // 资源ID
	Name   any // 资源名称
	Tips   any // 资源说明,用法提示
}
