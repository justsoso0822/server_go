// =================================================================================
// 代码由 GoFrame CLI 工具生成并维护。请勿编辑。
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// MemConfig 是表 mem_config 的 Go 结构体，用于 Where/Data 等 DAO 操作。
type MemConfig struct {
	g.Meta `orm:"table:mem_config, do:true"`
	Id     any //
	Value  any //
	Tips   any //
}
