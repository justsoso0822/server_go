// =================================================================================
// 代码由 GoFrame CLI 工具生成并维护。请勿编辑。
// =================================================================================

package entity

// UserData 是表 user_data 的 Go 结构体。
type UserData struct {
	Uid   int    `json:"uid"   orm:"uid"   description:""`
	Key   string `json:"key"   orm:"key"   description:""`
	Value string `json:"value" orm:"value" description:""`
}
