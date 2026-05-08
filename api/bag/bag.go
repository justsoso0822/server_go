package bag

import (
	"github.com/gogf/gf/v2/frame/g"
)

type GetBagReq struct {
	g.Meta  `path:"/bag/get_bag/:chapter" method:"get,post" tags:"Bag" summary:"获取用户背包"`
	Uid     int64 `json:"uid" v:"required"`
	Chapter int   `json:"chapter" in:"path" v:"required"`
}
type GetBagRes struct {
	g.Meta `mime:"application/json"`
}

type GetBagTpReq struct {
	g.Meta  `path:"/bag/get_bag_tp/:chapter" method:"get,post" tags:"Bag" summary:"获取用户背包tp"`
	Uid     int64 `json:"uid" v:"required"`
	Chapter int   `json:"chapter" in:"path" v:"required"`
}
type GetBagTpRes struct {
	g.Meta `mime:"application/json"`
}