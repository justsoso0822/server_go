package user

import (
	"server_go/internal/model"

	"github.com/gogf/gf/v2/frame/g"
)

type LoginReq struct {
	g.Meta   `path:"/user/login" method:"get,post" tags:"User" summary:"登录"`
	Uid      int64  `json:"uid" v:"required"`
	LoginKey string `json:"login_key"`
	Openid   string `json:"openid" v:"required"`
	Platform string `json:"platform" v:"required"`
	Version  string `json:"version" v:"required"`
}
type LoginRes model.LoginOutput

type AddTiliReq struct {
	g.Meta `path:"/user/add_tili" method:"get,post" tags:"User" summary:"测试增加体力"`
	Uid    int64 `json:"uid" v:"required"`
}
type AddTiliRes model.UpdateFieldOutput

type AddGoldReq struct {
	g.Meta `path:"/user/add_gold" method:"get,post" tags:"User" summary:"测试增加金币"`
	Uid    int64 `json:"uid" v:"required"`
}
type AddGoldRes model.UpdateFieldOutput

type AddDiamondReq struct {
	g.Meta `path:"/user/add_diamond" method:"get,post" tags:"User" summary:"测试增加钻石"`
	Uid    int64 `json:"uid" v:"required"`
}
type AddDiamondRes model.UpdateFieldOutput