// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// UserRes is the golang structure for table user_res.
type UserRes struct {
	Uid      int    `json:"uid"      orm:"uid"       description:""`
	Gold     int    `json:"gold"     orm:"gold"      description:""`
	Diamond  int    `json:"diamond"  orm:"diamond"   description:""`
	Star     int    `json:"star"     orm:"star"      description:""`
	Tili     int    `json:"tili"     orm:"tili"      description:""`
	TiliTime int    `json:"tiliTime" orm:"tili_time" description:""`
	Exp      int    `json:"exp"      orm:"exp"       description:""`
	Level    int    `json:"level"    orm:"level"     description:""`
	DayConf  string `json:"dayConf"  orm:"day_conf"  description:"每日重置的数据"`
	DayTime  int    `json:"dayTime"  orm:"day_time"  description:"上次重置时间"`
}
