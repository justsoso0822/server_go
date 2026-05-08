package game

import (
	"github.com/gogf/gf/v2/frame/g"
)

type OnlineReq struct {
	g.Meta  `path:"/game/online" method:"get,post" tags:"Game" summary:"记录在线时长"`
	Uid     int64 `json:"uid" v:"required"`
	Seconds int64 `json:"seconds" v:"required|min:0"`
}
type OnlineRes struct {
	Now int64 `json:"now"`
}

type TimeReq struct {
	g.Meta `path:"/game/time" method:"get,post" tags:"Game" summary:"获取服务器时间"`
}
type TimeRes struct {
	Now int64 `json:"now"`
}
