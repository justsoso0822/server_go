// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// PrfRes is the golang structure for table prf_res.
type PrfRes struct {
	Id   int    `json:"id"   orm:"id"   description:"资源ID"`
	Name string `json:"name" orm:"name" description:"资源名称"`
	Tips string `json:"tips" orm:"tips" description:"资源说明,用法提示"`
}
