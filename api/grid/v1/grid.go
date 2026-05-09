package v1

import (
	apiBag "server_go/api/bag/v1"

	"github.com/gogf/gf/v2/frame/g"
)

type GetGridReq struct {
	g.Meta  `path:"/grid/get/{chapter}" method:"get,post" tags:"Grid" summary:"获取棋盘数据"`
	Uid     int64 `json:"uid" v:"required"`
	Chapter int   `json:"chapter" in:"path" v:"required"`
}
type GetGridRes struct {
	Bag   *apiBag.BagRes `json:"bag,omitempty"`
	BagTp *apiBag.BagRes `json:"bag_tp,omitempty"`
	Tasks any            `json:"tasks,omitempty"`
}
