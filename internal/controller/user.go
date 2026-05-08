package controller

import (
	"context"

	apiUser "server_go/api/user"
	"server_go/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

var User = &cUser{}

type cUser struct{}

func (c *cUser) Login(ctx context.Context, req *apiUser.LoginReq) (res *apiUser.LoginRes, err error) {
	result, err := service.User().Login(ctx, req.Uid, req.LoginKey, req.Openid, req.Platform, req.Version)
	if err != nil {
		return nil, err
	}
	ghttp.RequestFromCtx(ctx).Response.WriteJson(result)
	return
}

func (c *cUser) AddTili(ctx context.Context, req *apiUser.AddTiliReq) (res *apiUser.AddTiliRes, err error) {
	uid := g.RequestFromCtx(ctx).Get("uid").Int64()
	ret := g.Map{}
	if err = service.User().UpdateTili(ctx, uid, 50, ret, "测试增加体力"); err != nil {
		return nil, err
	}
	ghttp.RequestFromCtx(ctx).Response.WriteJson(ret)
	return
}

func (c *cUser) AddGold(ctx context.Context, req *apiUser.AddGoldReq) (res *apiUser.AddGoldRes, err error) {
	uid := g.RequestFromCtx(ctx).Get("uid").Int64()
	ret := g.Map{}
	if err = service.User().UpdateGold(ctx, uid, 50, ret, "测试增加金币"); err != nil {
		return nil, err
	}
	ghttp.RequestFromCtx(ctx).Response.WriteJson(ret)
	return
}

func (c *cUser) AddDiamond(ctx context.Context, req *apiUser.AddDiamondReq) (res *apiUser.AddDiamondRes, err error) {
	uid := g.RequestFromCtx(ctx).Get("uid").Int64()
	ret := g.Map{}
	if err = service.User().UpdateDiamond(ctx, uid, 50, ret, "测试增加钻石"); err != nil {
		return nil, err
	}
	ghttp.RequestFromCtx(ctx).Response.WriteJson(ret)
	return
}