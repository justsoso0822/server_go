package user

import (
	"context"

	"server_go/api/user/v1"
	"server_go/internal/service"

	"github.com/gogf/gf/v2/database/gdb"
)

func (c *ControllerV1) Login(ctx context.Context, req *v1.LoginReq) (res *v1.LoginRes, err error) {
	out, err := service.User().Login(ctx, req.Uid, req.LoginKey, req.Openid, req.Platform, req.Version)
	if err != nil {
		return nil, err
	}
	return &v1.LoginRes{
		Uid:    out["uid"].(int64),
		Newbie: out["newbie"].(int),
		User:   out["user"],
		Res:    out["res"],
		Datas:  out["datas"].(gdb.Result),
		Items:  out["items"].(gdb.Result),
		Config: out["config"].(gdb.Result),
		Gm:     out["gm"].(int),
	}, nil
}
