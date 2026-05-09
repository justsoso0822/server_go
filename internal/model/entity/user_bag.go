// =================================================================================
// 代码由 GoFrame CLI 工具生成并维护。请勿编辑。
// =================================================================================

package entity

// UserBag 是表 user_bag 的 Go 结构体。
type UserBag struct {
	Id      int64  `json:"id"      orm:"id"      description:""`
	Uid     int    `json:"uid"     orm:"uid"     description:""`
	Chapter int    `json:"chapter" orm:"chapter" description:"副本id - 0-主格子"`
	Time    int    `json:"time"    orm:"time"    description:"操作更新时间"`
	Itemid  int    `json:"itemid"  orm:"itemid"  description:"物品id, 0=空"`
	Info    string `json:"info"    orm:"info"    description:"格子信息"`
	Type    int    `json:"type"    orm:"type"    description:"格子类型,0-普通,1-道具购买"`
}
