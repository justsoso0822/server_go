package user

import (
	"github.com/gogf/gf/v2/database/gdb"
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
type LoginRes struct {
	Uid    int64      `json:"uid"`
	Newbie int        `json:"newbie"`
	User   any        `json:"user"`
	Res    any        `json:"res,omitempty"`
	Datas  gdb.Result `json:"datas,omitempty"`
	Items  gdb.Result `json:"items,omitempty"`
	Config gdb.Result `json:"config,omitempty"`
	Gm     int        `json:"gm"`
}

type AddTiliReq struct {
	g.Meta `path:"/user/add_tili" method:"get,post" tags:"User" summary:"测试增加体力"`
	Uid    int64 `json:"uid" v:"required"`
}
type AddTiliRes UpdateFieldRes

type AddGoldReq struct {
	g.Meta `path:"/user/add_gold" method:"get,post" tags:"User" summary:"测试增加金币"`
	Uid    int64 `json:"uid" v:"required"`
}
type AddGoldRes UpdateFieldRes

type AddDiamondReq struct {
	g.Meta `path:"/user/add_diamond" method:"get,post" tags:"User" summary:"测试增加钻石"`
	Uid    int64 `json:"uid" v:"required"`
}
type AddDiamondRes UpdateFieldRes

type UpdateFieldRes struct {
	Res      any   `json:"res"`
	AddValue int64 `json:"add_value"`
}
