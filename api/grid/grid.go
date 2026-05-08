package grid

import (
	"server_go/internal/model"

	"github.com/gogf/gf/v2/frame/g"
)

type GetGridReq struct {
	g.Meta  `path:"/grid/get/:chapter" method:"get,post" tags:"Grid" summary:"获取棋盘数据"`
	Uid     int64 `json:"uid" v:"required"`
	Chapter int   `json:"chapter" in:"path" v:"required"`
}
type GetGridRes model.GridOutput