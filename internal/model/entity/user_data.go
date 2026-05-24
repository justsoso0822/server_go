// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// UserData is the golang structure for table user_data.
type UserData struct {
	Uid   int    `json:"uid"   orm:"uid"   description:""`
	Key   string `json:"key"   orm:"key"   description:""`
	Value string `json:"value" orm:"value" description:""`
}
