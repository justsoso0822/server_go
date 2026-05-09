package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

type ResVersionReq struct {
	g.Meta `path:"/res_version/{key}" method:"get,post" tags:"Other" summary:"获取资源版本号"`
	Key    string `json:"key" in:"path" v:"required"`
}
type ResVersionRes struct {
	Code int    `json:"code"`
	Ver  string `json:"ver,omitempty"`
	Msg  string `json:"msg,omitempty"`
}
