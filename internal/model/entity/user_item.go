// =================================================================================
// 代码由 GoFrame CLI 工具生成并维护。请勿编辑。
// =================================================================================

package entity

// UserItem 是表 user_item 的 Go 结构体。
type UserItem struct {
	Uid int `json:"uid" orm:"uid" description:""`
	Tid int `json:"tid" orm:"tid" description:""`
	Cnt int `json:"cnt" orm:"cnt" description:""`
}
