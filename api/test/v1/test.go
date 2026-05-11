package v1

import "github.com/gogf/gf/v2/frame/g"

type IndexReq struct {
	g.Meta `path:"/" method:"get,post" tags:"Test" summary:"测试接口"`
}
type IndexRes struct {
	Data any `json:"data"`
}

type TestDbReq struct {
	g.Meta `path:"db" method:"get,post" tags:"Test" summary:"测试数据库"`
	Uid    int `json:"uid" v:"required"`
}
type TestDbRes struct {
	Data any `json:"data"`
}
