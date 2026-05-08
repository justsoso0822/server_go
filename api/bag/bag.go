package bag

import (
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

type GetBagReq struct {
	g.Meta  `path:"/bag/get_bag/{chapter}" method:"get,post" tags:"Bag" summary:"获取用户背包"`
	Uid     int64 `json:"uid" v:"required"`
	Chapter int   `json:"chapter" in:"path" v:"required"`
}
type GetBagRes BagRes

type GetBagTpReq struct {
	g.Meta  `path:"/bag/get_bag_tp/{chapter}" method:"get,post" tags:"Bag" summary:"获取用户背包tp"`
	Uid     int64 `json:"uid" v:"required"`
	Chapter int   `json:"chapter" in:"path" v:"required"`
}
type GetBagTpRes BagRes

type BagRes struct {
	Uid     int64      `json:"uid"`
	Chapter int        `json:"chapter"`
	Bag     gdb.Result `json:"bag"`
}
