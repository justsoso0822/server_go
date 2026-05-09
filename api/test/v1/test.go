package v1

import "github.com/gogf/gf/v2/frame/g"

type IndexReq struct {
	g.Meta `path:"/" method:"get,post" tags:"Test" summary:"测试接口"`
}

type IndexRes struct {
	Data any `json:"data"`
}
