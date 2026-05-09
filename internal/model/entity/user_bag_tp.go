// =================================================================================
// 代码由 GoFrame CLI 工具生成并维护。请勿编辑。
// =================================================================================

package entity

// UserBagTp 是表 user_bag_tp 的 Go 结构体。
type UserBagTp struct {
	Id      int64 `json:"id"      orm:"id"      description:""`
	Uid     int   `json:"uid"     orm:"uid"     description:""`
	Chapter int   `json:"chapter" orm:"chapter" description:"副本id - 0 主线"`
	Time    int   `json:"time"    orm:"time"    description:"更新时间"`
	Itemid  int   `json:"itemid"  orm:"itemid"  description:"物品id, 0=空"`
	Count   int   `json:"count"   orm:"count"   description:"道具数量"`
}
