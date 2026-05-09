package model

import (
	"server_go/internal/model/entity"

	"github.com/gogf/gf/v2/database/gdb"
)

// --- 用户 ---

type LoginInput struct {
	Uid      int64  `json:"uid"`
	LoginKey string `json:"login_key"`
	Openid   string `json:"openid"`
	Platform string `json:"platform"`
	Version  string `json:"version"`
}

type LoginOutput struct {
	Uid    int64           `json:"uid"`
	Newbie int             `json:"newbie"`
	User   interface{}     `json:"user"`
	Res    *entity.UserRes `json:"res,omitempty"`
	Datas  gdb.Result      `json:"datas,omitempty"`
	Items  gdb.Result      `json:"items,omitempty"`
	Config gdb.Result      `json:"config,omitempty"`
	Gm     int             `json:"gm"`
}

type UpdateResInput struct {
	Uid    int64       `json:"uid"`
	Items  interface{} `json:"items"`
	Reason string      `json:"reason"`
}

type UpdateFieldInput struct {
	Uid    int64  `json:"uid"`
	Cnt    int64  `json:"cnt"`
	Reason string `json:"reason"`
}

type UpdateFieldOutput struct {
	Res      *entity.UserRes `json:"res"`
	AddValue int64           `json:"add_value"`
}

// --- 游戏 ---

type OnlineInput struct {
	Uid     int64 `json:"uid"`
	Seconds int64 `json:"seconds"`
}

// --- 背包 ---

type BagInput struct {
	Uid     int64 `json:"uid"`
	Chapter int   `json:"chapter"`
}

type BagOutput struct {
	Uid     int64      `json:"uid"`
	Chapter int        `json:"chapter"`
	Bag     gdb.Result `json:"bag"`
}

// --- 格子 ---

type GridOutput struct {
	Bag   *BagOutput  `json:"bag,omitempty"`
	BagTp *BagOutput  `json:"bag_tp,omitempty"`
	Tasks interface{} `json:"tasks,omitempty"`
}

// --- 其他 ---

type ResVersionOutput struct {
	Code int    `json:"code"`
	Ver  string `json:"ver,omitempty"`
	Msg  string `json:"msg,omitempty"`
}
