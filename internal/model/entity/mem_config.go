// =================================================================================
// 代码由 GoFrame CLI 工具生成并维护。请勿编辑。
// =================================================================================

package entity

// MemConfig 是表 mem_config 的 Go 结构体。
type MemConfig struct {
	Id    int    `json:"id"    orm:"id"    description:""`
	Value string `json:"value" orm:"value" description:""`
	Tips  string `json:"tips"  orm:"tips"  description:""`
}
