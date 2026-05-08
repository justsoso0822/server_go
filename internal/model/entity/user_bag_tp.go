// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// UserBagTp is the golang structure for table user_bag_tp.
type UserBagTp struct {
	Id      int64 `json:"id"      orm:"id"      description:""`
	Uid     int   `json:"uid"     orm:"uid"     description:""`
	Chapter int   `json:"chapter" orm:"chapter" description:"副本id - 0 主线"`
	Time    int   `json:"time"    orm:"time"    description:"更新时间"`
	Itemid  int   `json:"itemid"  orm:"itemid"  description:"物品id, 0=空"`
	Count   int   `json:"count"   orm:"count"   description:"道具数量"`
}
