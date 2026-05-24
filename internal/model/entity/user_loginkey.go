// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// UserLoginkey is the golang structure for table user_loginkey.
type UserLoginkey struct {
	Uid  int    `json:"uid"  orm:"uid"  description:""`
	Key  string `json:"key"  orm:"key"  description:""`
	Ver  int    `json:"ver"  orm:"ver"  description:"客户端版本号"`
	Time int    `json:"time" orm:"time" description:"登录时间"`
}
