package controller

import (
	"context"

	apiGame "server_go/api/game"
	"server_go/internal/service"

	"github.com/gogf/gf/v2/net/ghttp"
)

var Game = &cGame{}

type cGame struct{}

func (c *cGame) Online(ctx context.Context, req *apiGame.OnlineReq) (res *apiGame.OnlineRes, err error) {
	result, err := service.Game().Online(ctx, req.Uid, req.Seconds)
	if err != nil {
		return nil, err
	}
	ghttp.RequestFromCtx(ctx).Response.WriteJson(result)
	return
}

func (c *cGame) Time(ctx context.Context, req *apiGame.TimeReq) (res *apiGame.TimeRes, err error) {
	result := service.Game().ServerTime()
	ghttp.RequestFromCtx(ctx).Response.WriteJson(result)
	return
}