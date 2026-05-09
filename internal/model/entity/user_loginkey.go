// =================================================================================
// 代码由 GoFrame CLI 工具生成并维护。请勿编辑。
// =================================================================================

package entity

// UserLoginkey 是表 user_loginkey 的 Go 结构体。
type UserLoginkey struct {
	Uid  int    `json:"uid"  orm:"uid"  description:""`
	Key  string `json:"key"  orm:"key"  description:""`
	Ver  int    `json:"ver"  orm:"ver"  description:"客户端版本号"`
	Time int    `json:"time" orm:"time" description:"登录时间"`
}
