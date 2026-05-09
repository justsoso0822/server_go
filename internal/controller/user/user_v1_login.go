package user

import (
	"context"

	"server_go/api/user/v1"
	"server_go/internal/model"
	"server_go/internal/service"
)

func (c *ControllerV1) Login(ctx context.Context, req *v1.LoginReq) (res *v1.LoginRes, err error) {
	out, err := service.User().Login(ctx, &model.LoginInput{
		Uid:      req.Uid,
		LoginKey: req.LoginKey,
		Openid:   req.Openid,
		Platform: req.Platform,
		Version:  req.Version,
	})
	if err != nil {
		return nil, err
	}
	return &v1.LoginRes{
		Uid:    out.Uid,
		Newbie: out.Newbie,
		User:   out.User,
		Res:    out.Res,
		Datas:  out.Datas,
		Items:  out.Items,
		Config: out.Config,
		Gm:     out.Gm,
	}, nil
}
