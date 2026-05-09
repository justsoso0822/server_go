package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

type AddTiliReq struct {
	g.Meta `path:"/res/add_tili" method:"get,post" tags:"Res" summary:"测试增加体力"`
	Uid    int64 `json:"uid" v:"required"`
}
type AddTiliRes struct {
	Res     any   `json:"res"`
	AddTili int64 `json:"__add_tili"`
}

type AddGoldReq struct {
	g.Meta `path:"/res/add_gold" method:"get,post" tags:"Res" summary:"测试增加金币"`
	Uid    int64 `json:"uid" v:"required"`
}
type AddGoldRes struct {
	Res     any   `json:"res"`
	AddGold int64 `json:"__add_gold"`
}

type AddDiamondReq struct {
	g.Meta `path:"/res/add_diamond" method:"get,post" tags:"Res" summary:"测试增加钻石"`
	Uid    int64 `json:"uid" v:"required"`
}
type AddDiamondRes struct {
	Res        any   `json:"res"`
	AddDiamond int64 `json:"__add_diamond"`
}
